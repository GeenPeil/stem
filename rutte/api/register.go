package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/GeenPeil/stem/rutte/api/token"
	"github.com/GeenPeil/stem/rutte/common"
	"github.com/GeenPeil/stem/rutte/postcodenl"

	"github.com/pressly/chi/middleware"
	"github.com/pressly/chi/render"
	"github.com/rollick/decimal"
	mollieServices "github.com/rollick/gollie/services"
)

// registerStep1 creates a new account in the database, first step of member registration
func (a *API) registerStep1() http.HandlerFunc {
	log := a.log.WithField("handler", "register/step1")
	log.Infoln("setup")

	stmtPostMember, err := a.db.PrepareNamed(`
		INSERT INTO members.accounts
		(
			email,
			given_name,
			last_name,
			birthdate,
			phonenumber,
			registration_token
		) VALUES (
			:email,
			:given_name,
			:last_name,
			CAST(:birthdate AS date),
			:phonenumber,
			:registration_token
		) RETURNING id`)
	if err != nil {
		log.WithError(err).Fatal("error preparing statement")
	}

	type InMember struct {
		Email       string      `json:"email"`
		GivenName   string      `json:"givenName"`
		LastName    string      `json:"lastName"`
		Phonenumber string      `json:"phonenumber"`
		Birthdate   common.Date `json:"birthdate"`
	}

	type Member struct {
		InMember
		RegistrationToken string
	}

	type OutRegistrationToken struct {
		common.APIResponse
		RegistrationToken string `json:"registrationToken"`
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

		out := OutRegistrationToken{
			APIResponse: common.NewAPIResponse(),
		}
		if !common.RegexpValidateEmailAddress.MatchString(in.Email) {
			out.AddErrorCode(`rutte:invalid_email_address`)
			goto Done
		}

		{
			// member to save
			member := Member{
				InMember:          in,
				RegistrationToken: token.RandomToken(200),
			}
			err = stmtPostMember.QueryRow(member).Scan(&out.APIResponse.ID)
			if err != nil {
				if out.CheckPgErr(err) {
					goto Done
				}
				log.WithError(err).Error("error scanning row into struct")
				http.Error(w, "server error", http.StatusInternalServerError)
				return
			}
			out.RegistrationToken = member.RegistrationToken
		}
	Done:
		render.JSON(w, r, out)
	}
}

// registerStep2 sets address data on an account
func (a *API) registerStep2() http.HandlerFunc {
	log := a.log.WithField("handler", "register/step2")
	log.Infoln("setup")

	stmtUpdateAddress, err := a.db.PrepareNamed(`
		UPDATE members.accounts
			SET postalcode = :postalcode
			  , housenumber = CAST(:housenumber AS varchar)
			  , housenumber_suffix = :housenumber_suffix
			  , streetname = :streetname
			  , city = :city
			  , province = :province
			  , country_id = (SELECT id from i8n.countries WHERE code = :country_code)
		WHERE id = :id AND registration_token = :registration_token`)
	if err != nil {
		log.WithError(err).Fatal("error preparing statement")
	}

	type InAddress struct {
		ID                uint64 `json:"id"`
		RegistrationToken string `json:"registrationToken"`

		Postalcode        string `json:"postalcode"`
		Housenumber       uint64 `json:"housenumber"`
		HousenumberSuffix string `json:"housenumberSuffix"`
		Streetname        string `json:"streetname"`
		City              string `json:"city"`
		Province          string `json:"province"`
		CountryCode       string `json:"countryCode"`
	}

	type OutAddressSaved struct {
		common.APIResponse
		Streetname string `json:"streetname,omitempty"`
		City       string `json:"city,omitempty"`
		Province   string `json:"province,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		log := log.WithField("requestID", middleware.GetReqID(r.Context()))

		// read profile from db
		var in = InAddress{}
		err = json.NewDecoder(r.Body).Decode(&in)
		r.Body.Close()
		if err != nil {
			log.WithError(err).Warn("invalid json")
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		out := OutAddressSaved{
			APIResponse: common.NewAPIResponse(),
		}
		{
			if in.CountryCode == `NL` {
				postcodeData, err := a.postcodeAPI.Check(in.Postalcode, in.Housenumber)
				if err != nil {
					if exception, ok := err.(postcodenl.Exception); ok {
						out.AddErrorCode(`rutte:` + exception.ExceptionID)
						goto Done
					}
					if err == postcodenl.ErrInvalidPostcode {
						out.AddErrorCode("rutte:invalid_postal_code")
						goto Done
					}
					log.WithError(err).Error("error doing postcode check")
					http.Error(w, "server error", http.StatusInternalServerError)
					return
				}
				in.Streetname = postcodeData.Street
				in.City = postcodeData.City
				in.Province = postcodeData.Province
			}

			res, err := stmtUpdateAddress.Exec(in)
			if err != nil {
				if out.CheckPgErr(err) {
					goto Done
				}
				log.WithError(err).Error("error scanning row into struct")
				http.Error(w, "server error", http.StatusInternalServerError)
				return
			}
			affected, err := res.RowsAffected()
			if err != nil {
				log.WithError(err).Error("error getting affected rows")
				http.Error(w, "server error", http.StatusInternalServerError)
				return
			}
			if affected == 0 {
				// perhaps the registration token wasn't valid anymore?
				log.Warn("account was not modified")
				out.AddErrorCode("rutte:account_not_modified")
				return
			}

			if in.CountryCode == `NL` {
				// push saved streetname, city, province back to frontend
				out.Streetname = in.Streetname
				out.City = in.City
				out.Province = in.Province
			}
		}
	Done:
		render.JSON(w, r, &out)
	}
}

// registerStep3 creates a new mollie payment
func (a *API) registerStep3() http.HandlerFunc {
	log := a.log.WithField("handler", "register/step3")
	log.Infoln("setup")

	stmtCreatePayment, err := a.db.PrepareNamed(`
		INSERT INTO members.payments (
			account_id,
			token,
			mollie_id,
			mollie_status,
			mollie_created
		) VALUES (
			(SELECT id FROM members.accounts WHERE id = :id AND registration_token = :registration_token),
			:token,
			:mollie_id,
			CAST(:mollie_status AS members.enum_mollie_status),
			:mollie_created
		)`)
	if err != nil {
		log.WithError(err).Fatalln("error preparing statement")
	}

	type InRequest struct {
		ID                uint64 `json:"id"`
		RegistrationToken string `json:"registrationToken"`
	}

	type Payment struct {
		InRequest
		Token         string
		MollieID      string
		MollieStatus  string
		MollieCreated *time.Time
	}

	type OutResponse struct {
		common.APIResponse
		PaymentURL string `json:"paymentURL"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		log := log.WithField("requestID", middleware.GetReqID(r.Context()))

		// read profile from db
		var in = InRequest{}
		err := json.NewDecoder(r.Body).Decode(&in)
		r.Body.Close()
		if err != nil {
			log.WithError(err).Warn("invalid json")
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		out := OutResponse{
			APIResponse: common.NewAPIResponse(),
		}
		{
			var paymentToken = token.RandomToken(80)
			molliePayment, _, err := a.molliePaymentService.Create(&mollieServices.PaymentRequest{
				Amount:      decimal.New(12, 0),
				Description: `Lidmaatschap GeenPeil 2017`,
				RedirectUrl: a.selfHTTPAddress + `/lid-worden/check-payment?token=` + paymentToken,
				WebhookUrl:  a.selfHTTPAddress + `/api/register/mollie-webhook`,
				Locale:      "nl_NL",
				// Metadata:    map[string]string{},
			})
			if err != nil {
				log.WithError(err).Errorln(`error creating mollie payment`)
			}

			payment := Payment{
				InRequest:     in,
				Token:         paymentToken,
				MollieID:      molliePayment.ID,
				MollieStatus:  molliePayment.Status,
				MollieCreated: molliePayment.CreatedDatetime,
			}
			res, err := stmtCreatePayment.Exec(payment)
			if err != nil {
				if out.CheckPgErr(err) {
					goto Done
				}
				log.WithError(err).Error("error scanning row into struct")
				http.Error(w, "server error", http.StatusInternalServerError)
				return
			}
			affected, err := res.RowsAffected()
			if err != nil {
				log.WithError(err).Error("error getting affected rows")
				http.Error(w, "server error", http.StatusInternalServerError)
				return
			}
			if affected == 0 {
				// perhaps the registration token wasn't valid anymore?
				log.Warn("payment was not inserted")
				out.AddErrorCode("rutte:payment_not_created")
				return
			}
			out.PaymentURL = molliePayment.Links.PaymentUrl
		}
	Done:
		render.JSON(w, r, &out)
	}
}

// paymentStatus
func (a *API) mollieCheckPayment() http.HandlerFunc {
	log := a.log.WithField("handler", "register/check-payment")
	log.Infoln("setup")

	stmtSelectMollieID, err := a.db.Preparex(`SELECT mollie_id FROM members.payments WHERE token = $1`)
	if err != nil {
		log.WithError(err).Fatalln("error preparing statement")
	}

	type OutPaymentStatus struct {
		common.APIResponse
		Paid bool `json:"paid"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		log := log.WithField("requestID", middleware.GetReqID(r.Context()))

		out := OutPaymentStatus{
			APIResponse: common.NewAPIResponse(),
		}
		{
			token := r.FormValue("token")
			var mollieID string
			err := stmtSelectMollieID.QueryRowx(token).Scan(&mollieID)
			if err != nil {
				if err == sql.ErrNoRows {
					out.AddErrorCode(`rutte:invalid_payment_token`)
					goto Done
				}
				log.WithError(err).Errorln("error getting mollieID by token")
			}
			molliePayment, _, err := a.molliePaymentService.Fetch(mollieID)
			if err != nil {
				log.WithError(err).Errorln(`error getting mollie payment`)
				http.Error(w, "server error", http.StatusInternalServerError)
				return
			}
			if molliePayment.PaidDatetime != nil {
				out.Paid = true
			}
		}
	Done:
		render.JSON(w, r, &out)
	}
}

// mollieWebhook handles mollie calls
func (a *API) mollieWebhook() http.HandlerFunc {
	log := a.log.WithField("handler", "register/mollie-webhook")
	log.Infoln("setup")

	stmtUpdatePayment, err := a.db.PrepareNamed(`
		UPDATE members.payments
			SET mollie_status = CAST(:mollie_status AS members.enum_mollie_status)
			  , mollie_paid = :mollie_paid
		WHERE mollie_id = :mollie_id
		`)
	if err != nil {
		log.WithError(err).Fatalln("error preparing statement")
	}

	type Payment struct {
		MollieID     string
		MollieStatus string
		MolliePaid   *time.Time
	}

	return func(w http.ResponseWriter, r *http.Request) {
		log := log.WithField("requestID", middleware.GetReqID(r.Context()))

		mollieID := r.FormValue("id")
		molliePayment, _, err := a.molliePaymentService.Fetch(mollieID)
		if err != nil {
			log.WithError(err).Errorln(`error getting mollie payment`)
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		payment := Payment{
			MollieID:     molliePayment.ID,
			MollieStatus: molliePayment.Status,
			MolliePaid:   molliePayment.PaidDatetime,
		}
		res, err := stmtUpdatePayment.Exec(payment)
		if err != nil {
			log.WithError(err).Error("error executing query")
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		affected, err := res.RowsAffected()
		if err != nil {
			log.WithError(err).Error("error getting affected rows")
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		if affected == 0 {
			log.Warn("payment was not updated")
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

	}
}
