package api

import (
	"net/http"

	"github.com/pressly/chi/middleware"
	"github.com/pressly/chi/render"
)

// getCountryList retrieves the profile overview for currently logged in user
func (a *API) getCountryList() http.HandlerFunc {
	log := a.log.WithField("handler", "country-list")
	log.Infoln("setup")

	stmtSelectCountries, err := a.db.Preparex(`
		SELECT countries.code, country_names.name
		FROM i8n.countries
			INNER JOIN i8n.country_names
				ON countries.id = country_names.country_id
		WHERE country_names.language_id = (SELECT id FROM i8n.languages WHERE code = $1)`)
	if err != nil {
		log.WithError(err).Fatal("error preparing stmtSelectCountries")
	}

	type OutCountry struct {
		Code string `json:"code"`
		Name string `json:"name"`
	}

	type OutCountryList struct {
		Countries []OutCountry `json:"countries"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		log := log.WithField("requestID", middleware.GetReqID(r.Context()))

		languageCode := `nl_NL`
		out := OutCountryList{}

		rows, err := stmtSelectCountries.Queryx(languageCode)
		if err != nil {
			log.WithError(err).Error("error getting countries from db")
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		for rows.Next() {
			var c OutCountry
			err = rows.StructScan(&c)
			if err != nil {
				log.WithError(err).Error("error scanning country")
				http.Error(w, "server error", http.StatusInternalServerError)
				return
			}
			out.Countries = append(out.Countries, c)
		}

		render.JSON(w, r, &out)
	}
}
