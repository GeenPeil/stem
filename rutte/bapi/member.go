package bapi

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/pressly/chi/render"
)

// regexpValidateEmailAddress is just a very simple regular expression to catch 99% of user input mistakes.
// Email addresses in the system are always validated by actually sending a verification email to it.
var regexpValidateEmailAddress = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// postMember a new member into the database
func (a *API) postMember() http.HandlerFunc {
	log := a.log.WithField("handler", "putMember")
	log.Infoln("setup")

	stmtPostMember, err := a.db.PrepareNamed(`
		INSERT INTO members.accounts
		(
			email,
			nickname,
			given_name,
			first_names,
			initials,
			last_name,
			birthdate,
			phonenumber,
			postalcode,
			housenumber,
			housenumber_suffix,
			streetname,
			city,
			province,
			country
		) VALUES (
			:email,
			:nickname,
			:given_name,
			:first_names,
			:initials,
			:last_name,
			CAST(:birthdate AS date),
			:phonenumber,
			:postalcode,
			:housenumber,
			:housenumber_suffix,
			:streetname,
			:city,
			:province,
			:country
		) RETURNING id`)
	if err != nil {
		log.WithError(err).Fatal("error preparing statement")
	}

	type InMember struct {
		Email             string  `json:"email"`
		Nickname          *string `json:"nickname"`
		GivenName         string  `json:"givenName"`
		FirstNames        string  `json:"firstNames"`
		Initials          string  `json:"initials"`
		LastName          string  `json:"lastName"`
		Birthdate         Date    `json:"birthdate"`
		Phonenumber       string  `json:"phonenumber"`
		Postalcode        string  `json:"postalcode"`
		Housenumber       string  `json:"housenumber"`
		HousenumberSuffix string  `json:"housenumberSuffix"`
		Streetname        string  `json:"streetname"`
		City              string  `json:"city"`
		Province          string  `json:"province"`
		Country           string  `json:"country"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		log := log.WithField("requestID", middleware.GetReqID(r.Context()))

		// read profile from db
		var in = InMember{}
		err = json.NewDecoder(r.Body).Decode(&in)
		r.Body.Close()
		if err != nil {
			log.WithError(err).Warn("invalid json")
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		response := NewAPIResponse()
		if !regexpValidateEmailAddress.MatchString(in.Email) {
			response.AddErrorCode(`rutte:invalid_email_address`)
			goto Done
		}

		err = stmtPostMember.QueryRow(in).Scan(&response.ID)
		if err != nil && !response.CheckPgErr(err) {
			log.WithError(err).Error("error scanning row into struct")
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

	Done:
		render.JSON(w, r, response)
	}
}

// getMember retrieves the member details for given account ID
func (a *API) getMember() http.HandlerFunc {
	log := a.log.WithField("handler", "getMember")
	log.Infoln("setup")

	stmtGetMember, err := a.db.Preparex(`
		SELECT
			accounts.id,
			accounts.email,
			accounts.nickname,
			accounts.given_name,
			accounts.first_names,
			accounts.initials,
			accounts.last_name,
			to_char(accounts.birthdate, 'YYYY-MM-DD') AS birthdate,
			accounts.phonenumber,
			accounts.postalcode,
			accounts.housenumber,
			accounts.housenumber_suffix,
			accounts.streetname,
			accounts.city,
			accounts.province,
			accounts.country,
			to_char(accounts.fee_last_payment_date, 'YYYY-MM-DD') AS fee_last_payment_date,
			accounts.fee_paid,
			accounts.is_adult,
			accounts.verified_email,
			accounts.verified_identity,
			accounts.verified_voting_entitlement
		FROM members.accounts
		WHERE accounts.id = $1`)
	if err != nil {
		log.WithError(err).Fatal("error preparing statement")
	}

	type OutMemberDetails struct {
		ID                        uint64  `json:"id"`
		Email                     string  `json:"email"`
		Nickname                  *string `json:"nickname"`
		GivenName                 string  `json:"givenName"`
		FirstNames                string  `json:"firstNames"`
		Initials                  string  `json:"initials"`
		LastName                  string  `json:"lastName"`
		Birthdate                 Date    `json:"birthdate"`
		IsAdult                   bool    `json:"isAdult"`
		Phonenumber               string  `json:"phonenumber"`
		Postalcode                string  `json:"postalcode"`
		Housenumber               string  `json:"housenumber"`
		HousenumberSuffix         string  `json:"housenumberSuffix"`
		Streetname                string  `json:"streetname"`
		City                      string  `json:"city"`
		Province                  string  `json:"province"`
		Country                   string  `json:"country"`
		FeeLastPaymentDate        *Date   `json:"feeLastPaymentDate"`
		FeePaid                   bool    `json:"feePaid"`
		VerifiedEmail             bool    `json:"verifiedEmail"`
		VerifiedIdentity          bool    `json:"verifiedIdentity"`
		VerifiedVotingEntitlement bool    `json:"verifiedVotingEntitlement"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		log := log.WithField("requestID", middleware.GetReqID(r.Context()))

		accountID, err := strconv.ParseUint(chi.URLParam(r, `id`), 10, 64)
		if err != nil {
			log.WithError(err).Warn("invalid account ID provided")
			http.Error(w, "invalid accountID", http.StatusBadRequest)
			return
		}

		// read profile from db
		var out = OutMemberDetails{}
		err = stmtGetMember.QueryRowx(accountID).StructScan(&out)
		if err != nil {
			log.WithError(err).Error("error scanning row into struct")
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		render.JSON(w, r, out)
	}
}

// putMember updates an existing member.
func (a *API) putMember() http.HandlerFunc {
	log := a.log.WithField("handler", "putMember")
	log.Infoln("setup")

	stmtPutMember, err := a.db.PrepareNamed(`
		UPDATE members.accounts
		SET
			nickname = :nickname,
			given_name = :given_name,
			first_names = :first_names,
			initials = :initials,
			last_name = :last_name,
			birthdate = CAST(:birthdate AS date),
			phonenumber = :phonenumber,
			postalcode = :postalcode,
			housenumber = :housenumber,
			housenumber_suffix = :housenumber_suffix,
			streetname = :streetname,
			city = :city,
			province = :province,
			country = :country
		WHERE id = :id`)
	if err != nil {
		log.WithError(err).Fatal("error preparing statement")
	}

	type InMember struct {
		ID                uint64  `json:"id"`
		Nickname          *string `json:"nickname"`
		GivenName         string  `json:"givenName"`
		FirstNames        string  `json:"firstNames"`
		Initials          string  `json:"initials"`
		LastName          string  `json:"lastName"`
		Birthdate         Date    `json:"birthdate"`
		Phonenumber       string  `json:"phonenumber"`
		Postalcode        string  `json:"postalcode"`
		Housenumber       string  `json:"housenumber"`
		HousenumberSuffix string  `json:"housenumberSuffix"`
		Streetname        string  `json:"streetname"`
		City              string  `json:"city"`
		Province          string  `json:"province"`
		Country           string  `json:"country"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		log := log.WithField("requestID", middleware.GetReqID(r.Context()))

		accountID, err := strconv.ParseUint(chi.URLParam(r, `id`), 10, 64)
		if err != nil {
			log.WithError(err).Warn("invalid account ID provided")
			http.Error(w, "invalid accountID", http.StatusBadRequest)
			return
		}

		// read profile from db
		var in = InMember{}
		err = json.NewDecoder(r.Body).Decode(&in)
		r.Body.Close()
		if err != nil {
			log.WithError(err).Warn("invalid json")
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		if in.ID != accountID {
			log.WithError(err).Warn("query accountID does not match post data account ID")
			http.Error(w, "ID mismatch", http.StatusBadRequest)
			return

		}

		response := NewAPIResponse()
		_, err = stmtPutMember.Exec(in)
		if err != nil && !response.CheckPgErr(err) {
			log.WithError(err).Error("error scanning row into struct")
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		render.JSON(w, r, response)
	}
}
