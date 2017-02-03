package api

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/jpeg" // Allow parsing of JPEG images
	_ "image/png"  // Allow parsing of PNG images
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/pressly/chi/render"

	"rsc.io/pdf"
)

func (a *API) getEntitlementsOverview() http.HandlerFunc {
	log := a.log.WithField("handler", "entitlement/overview")
	log.Infoln("setup")

	stmt, err := a.db.Preparex(`select  id
		        									,       created_at
															, 			original_filename
		                          from    members.entitlement_proofs
															where   account_id = $1`)
	if err != nil {
		log.WithError(err).Fatal("error preparing statement")
	}

	type OutEntitlement struct {
		ID               uint64    `json:"id"`
		CreatedAt        time.Time `json:"created_at"`
		OriginalFilename string    `json:"original_filename"`
	}

	type Out struct {
		Entitlements []OutEntitlement `json:"entitlements"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var out = Out{}

		accountID, err := a.accountID(r)
		if err != nil {
			log.WithError(err).Error("error retrieving account identifier")
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		rows, err := stmt.Queryx(accountID)
		if err != nil {
			log.WithError(err).Error("error retrieving entitlements from db")
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		for rows.Next() {
			var e OutEntitlement
			err = rows.StructScan(&e)

			if err != nil {
				log.WithError(err).Error("error scanning entitlement")
				http.Error(w, "server error", http.StatusInternalServerError)
				return
			}

			out.Entitlements = append(out.Entitlements, e)
		}

		render.JSON(w, r, &out)
	}
}

func (a *API) createEntitlement() http.HandlerFunc {
	log := a.log.WithField("handler", "entitlement/create")
	log.Infoln("setup")

	stmt, err := a.db.Preparex(`INSERT INTO members.entitlement_proofs (account_id, filename, original_filename)
														  VALUES ($1, $2, $3)
															RETURNING entitlement_proofs.id`)
	if err != nil {
		log.WithError(err).Fatal("error preparing statement")
	}

	type OutEntitlement struct {
		ID               uint64 `json:"id"`
		OriginalFileName string `json:"original_filename"`
	}

	type Out struct {
		Entitlement OutEntitlement `json:"entitlement"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		const limit = (1 << 20) * 10 // 10MB
		r.Body = http.MaxBytesReader(w, r.Body, limit)

		file, handle, err := r.FormFile("upload")
		if err != nil {
			log.WithError(err).Error("error uploading file")
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		originalFileName := handle.Filename
		if originalFileName == "" || len(originalFileName) > 200 {
			log.WithError(err).Error("error due to invalid file name")
			http.Error(w, "invalid file name", http.StatusForbidden)
			return
		}

		f, err := os.OpenFile(handle.Filename, os.O_RDONLY, 0644)
		if err != nil {
			log.WithError(err).Error("error creating readonly file")
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		defer f.Close()

		b, err := ioutil.ReadAll(f)
		if err != nil {
			log.WithError(err).Error("error reading file contents")
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		extension, err := validateFileType(b, f)
		if err != nil {
			log.WithError(err).Error("invalid file type")
			http.Error(w, "invalid file type", http.StatusForbidden)
			return
		}

		// Now save to the file system. The following considerations apply:
		// 1. The file must have a unique file name, which consists of a pseudo-random
		// sequence of base64 encoded bytes, current UTC time and file extension, i.e.:
		// abcdefghijklmnopqrstuvwxyz1234567890-YYYYMMDDHHMISS.pdf
		// 3. The values gathered above are returned in a JSON object.
		s, err := randomString(32)
		if err != nil {
			log.WithError(err).Error("error creating random string")
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		t := time.Now().UTC()
		now := fmt.Sprintf("%d%02d%02d%02d%02d%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

		fileName := fmt.Sprintf("%s-%s%s", s, now, extension)

		// Now save the file to filesystem.
		const uploadsFolder = "uploads/"
		err = ioutil.WriteFile(fmt.Sprintf("%s%s", uploadsFolder, fileName), b, 0644)
		if err != nil {
			log.WithError(err).Error("error saving file to file system")
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		accountID, err := a.accountID(r)
		if err != nil {
			log.WithError(err).Error("error retrieving account identifier")
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		var out = Out{}

		err = stmt.QueryRow(accountID, fileName, originalFileName).Scan(&out.Entitlement.ID)
		if err != nil {
			log.WithError(err).Error("error inserting entitlement metadata into database")
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		out.Entitlement.OriginalFileName = originalFileName

		render.JSON(w, r, &out)
	}
}

func validateFileType(b []byte, f *os.File) (string, error) {
	fileType := http.DetectContentType(b)

	switch fileType {
	case "application/pdf":
		{
			// Additional checks to verify the integrity of the PDF file.
			fi, err := f.Stat()
			if err != nil {
				return "", err
			}

			if _, err := pdf.NewReader(f, fi.Size()); err != nil {
				return "", err
			}

			return ".pdf", nil
		}
	case "image/png":
		{
			// Attempt to parse the PNG image to verify its integrity.
			_, _, err := image.Decode(f)
			if err != nil {
				return "", err
			}

			return ".png", nil
		}
	case "image/jpeg":
		{
			// Attempt to parse the JPEG image to verify its integrity.
			_, _, err := image.Decode(f)
			if err != nil {
				return "", err
			}

			return ".jpg", nil
		}
	}

	return "", fmt.Errorf("file type '%s' is not supported", fileType)
}

func randomString(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)

	return base64.URLEncoding.EncodeToString(b), err
}
