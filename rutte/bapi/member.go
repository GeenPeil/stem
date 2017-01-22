package bapi

import (
	"net/http"
	"strconv"

	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/pressly/chi/render"
)

// getProfile retrieves the profile overview for currently logged in user
func (a *API) getMember() http.HandlerFunc {
	log := a.log.WithField("handler", "member")
	log.Infoln("setup")

	stmtGetItem, err := a.db.Preparex(`
		SELECT
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
		FeeLastPaymentDate        Date    `json:"feeLastPaymentDate"`
		FeePaid                   bool    `json:"feePaid"`
		VerifiedEmail             bool    `json:"verifiedEmail"`
		VerifiedIdentity          bool    `json:"verifiedIdentity"`
		VerifiedVotingEntitlement bool    `json:"verifiedVotingEntitlement"`
	}

	type Out struct {
		MemberDetails OutMemberDetails `json:"memberDetails"`
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
		var out = Out{
			MemberDetails: OutMemberDetails{},
		}
		err = stmtGetItem.QueryRowx(accountID).StructScan(&out.MemberDetails)
		if err != nil {
			log.WithError(err).Error("error scanning row into struct")
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		render.JSON(w, r, out)
	}
}
