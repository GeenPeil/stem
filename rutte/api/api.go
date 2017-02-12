package api

import (
	"net/http"

	"github.com/GeenPeil/stem/rutte/postcodenl"
	"github.com/GeenPeil/stem/rutte/version"

	"github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/pressly/chi/render"
	mollieServices "github.com/rollick/gollie/services"
)

// API service for cpserver
type API struct {
	log                  *logrus.Entry
	db                   *sqlx.DB
	postcodeAPI          postcodenl.API
	molliePaymentService *mollieServices.PaymentService
	selfHTTPAddress      string
}

// New creates a new instance of the API service
func New(log *logrus.Entry, db *sqlx.DB, postcodeAPI postcodenl.API, molliePaymentService *mollieServices.PaymentService, selfHTTPAddress string) *API {
	return &API{
		log:                  log,
		db:                   db,
		postcodeAPI:          postcodeAPI,
		molliePaymentService: molliePaymentService,
		selfHTTPAddress:      selfHTTPAddress,
	}
}

// AttachChiRouter attaches API routes to the provided chi router.
func (a *API) AttachChiRouter(r chi.Router) {
	a.log.Infoln("attach chi router api")
	r.Use(middleware.NoCache)
	r.Get("/", a.getRoot)
	r.Get("/country-list", a.getCountryList())
	r.Post("/register/step1", a.registerStep1())
	r.Post("/register/step2", a.registerStep2())
	r.Post("/register/step3", a.registerStep3())
	r.Get("/register/check-payment", a.mollieCheckPayment())
	r.Post("/register/mollie-webhook", a.mollieWebhook())
	r.Post("/login", a.login())
	r.Post("/logout", a.logout())
	r.Route("/private", func(r chi.Router) {
		r.Use(a.newAccountAuthMiddleware())
		r.Get("/profile", a.getProfileOverview())
	})
}

// getRoot handles get requests to root (show api version?)
func (a *API) getRoot(w http.ResponseWriter, r *http.Request) {
	out := map[string]string{
		"apiName":        "rutte-api",
		"apiDescription": "API's for voting platform",
		"apiVersion":     "v0.0",
		"appVersion":     version.String(),
	}
	render.JSON(w, r, out)
}
