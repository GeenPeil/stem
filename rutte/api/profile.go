package api

import (
	"net/http"

	"github.com/pressly/chi/middleware"
	"github.com/pressly/chi/render"
)

// getProfile retrieves the profile overview for currently logged in user
func (a *API) getProfileOverview() http.HandlerFunc {
	log := a.log.WithField("handler", "profile/overview")
	log.Infoln("setup")

	stmtGetItem, err := a.db.Preparex(`
		SELECT
			accounts.id,
			accounts.nickname,
			accounts.email
		FROM members.accounts
		WHERE accounts.id = $1`)
	if err != nil {
		log.WithError(err).Fatal("error preparing statement")
	}

	type OutProfile struct {
		ID       uint64 `json:"id"`
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		log := log.WithField("requestID", middleware.GetReqID(r.Context()))

		// obtain accountID from authenticated session
		accountID, err := a.accountID(r)
		if err != nil {
			log.WithError(err).Error("could not get account id")
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		// read profile from db
		var out = OutProfile{}
		err = stmtGetItem.QueryRowx(accountID).StructScan(&out)
		if err != nil {
			log.WithError(err).Error("error scanning row into struct")
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		render.JSON(w, r, out)
	}
}
