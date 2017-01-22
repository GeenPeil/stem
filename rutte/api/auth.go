package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/GeenPeil/stem/rutte/api/token"
	"github.com/lib/pq"

	"github.com/pressly/chi/middleware"
	"github.com/pressly/chi/render"
)

// ErrInvalidAccountAuth is returned by (a *API).accountID(…) when the http.Request does
// not cary valid sessiontoken.
var ErrInvalidAccountAuth = errors.New("invalid account authentication")

var headerSessionToken = `x-gp-sessiontoken`

type ctxKey string

var ctxKeyAccountAuth = ctxKey("account-auth")

// (a *API) newAccountAuthMiddleware returns a middleware handler that validates authentication and
// sets up a account authentication object in the request context.
func (a *API) newAccountAuthMiddleware() func(next http.Handler) http.Handler {
	stmtObtainSession, err := a.db.Preparex(`SELECT sessions.account_id FROM members.sessions where sessions.token = $1`)
	if err != nil {
		a.log.WithError(err).Fatalln("error preparing statement to obtain token")
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: validate authentication
			token := r.Header.Get(headerSessionToken)
			if len(token) == 0 {
				http.Error(w, "invalid token", http.StatusForbidden)
				return
			}
			var accountID uint64
			err := stmtObtainSession.QueryRow(token).Scan(&accountID)
			if err != nil {
				a.log.WithError(err).Errorln("error obtaining session")
				http.Error(w, "error obtaining session", http.StatusInternalServerError)
			}

			// Add authentication to context
			ctx := context.WithValue(r.Context(), ctxKeyAccountAuth, accountID)

			// Call the next item in the chain
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// accountID returns the accountID or an error for the given http.Request.
// http.Request must've ran through the PartnerAuthMiddleware to be valid.
func (a *API) accountID(r *http.Request) (uint64, error) {
	ctx := r.Context()
	accountID, ok := ctx.Value(ctxKeyAccountAuth).(uint64)
	if !ok {
		return 0, ErrInvalidAccountAuth
	}
	return accountID, nil
}

func (a *API) login() http.HandlerFunc {
	log := a.log.WithField("handler", "login")
	log.Infoln("setup")

	stmtNewSession, err := a.db.Preparex(`INSERT INTO members.sessions (account_id, token) values ($1, $2) RETURNING sessions.id`)
	if err != nil {
		log.WithError(err).Fatal("error preparing statement")
	}

	type InLoginRequest struct {
		ID uint64 `json:"id"`
	}

	type OutSession struct {
		ID    uint64 `json:"id"`
		Token string `json:"token"`
	}

	type Out struct {
		Session OutSession `json:"session"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		log := log.WithField("requestID", middleware.GetReqID(r.Context()))

		inLoginRequest := InLoginRequest{}
		defer r.Body.Close()
		err = json.NewDecoder(r.Body).Decode(&inLoginRequest)
		if err != nil {
			log.WithError(err).Error("error decoding json")
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		var out = Out{}
		var token = token.RandomToken(200)
		err := stmtNewSession.QueryRow(&inLoginRequest.ID, token).Scan(&out.Session.ID)
		if err != nil {
			if perr, ok := err.(*pq.Error); ok && perr.Code == `23503` && perr.Constraint == `sessions_fk_account` {
				http.Error(w, "invalid account ID", http.StatusForbidden)
				return
			}
			log.WithError(err).Error("error creating new session")
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		out.Session.Token = token
		render.JSON(w, r, &out)
	}
}
