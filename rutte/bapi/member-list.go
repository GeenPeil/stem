package bapi

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/GeenPeil/stem/rutte/common"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pressly/chi/middleware"
	"github.com/pressly/chi/render"
)

// getMember retrieves the member details for given account ID
func (a *API) getMemberList() http.HandlerFunc {
	log := a.log.WithField("handler", "getMemberList")
	log.Infoln("setup")

	regexpQueryReplacer := regexp.MustCompile(`([^&|])\s+([^&|])`)

	const queryStart = `
SELECT
	accounts.id,
	accounts.given_name,
	accounts.nickname,
	accounts.initials,
	accounts.last_name
FROM members.accounts
`

	const queryEnd = `
LIMIT 100
`

	stmtMemberListByID, err := a.db.Preparex(queryStart + `WHERE accounts.id = CAST($1 AS integer)` + queryEnd)
	if err != nil {
		log.WithError(err).Fatal("error preparing statement")
	}

	stmtMemberListByNickname, err := a.db.Preparex(queryStart + `WHERE accounts.nickname = $1` + queryEnd)
	if err != nil {
		log.WithError(err).Fatal("error preparing statement")
	}

	stmtMemberListByEmail, err := a.db.Preparex(queryStart + `WHERE accounts.email = $1` + queryEnd)
	if err != nil {
		log.WithError(err).Fatal("error preparing statement")
	}

	stmtMemberListByPostalcode, err := a.db.Preparex(queryStart + `WHERE accounts.postalcode = $1 ORDER BY housenumber` + queryEnd)
	if err != nil {
		log.WithError(err).Fatal("error preparing statement")
	}

	stmtMemberListByTextsearch, err := a.db.Preparex(queryStart + `WHERE accounts.textsearch_vector @@ to_tsquery('english', $1)` + queryEnd)
	if err != nil {
		log.WithError(err).Fatal("error preparing statement")
	}

	type OutMember struct {
		ID        uint64  `json:"id"`
		GivenName string  `json:"givenName"`
		Nickname  *string `json:"nickname"`
		Initials  string  `json:"initials"`
		LastName  string  `json:"lastName"`
	}

	type OutMemberList struct {
		Error   string      `json:"error"`
		Members []OutMember `json:"members"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		log := log.WithField("requestID", middleware.GetReqID(r.Context()))

		searchQuery := r.FormValue(`searchQuery`)
		if len(searchQuery) == 0 {
			http.Error(w, "empty searchQuery", http.StatusBadRequest)
			log.Warn("empty searchQuery")
			return
		}

		var queryStatement *sqlx.Stmt
		if searchQuery[0] == '@' {
			// search by nickname
			queryStatement = stmtMemberListByNickname
			searchQuery = searchQuery[1:]
		} else if _, err := strconv.ParseUint(searchQuery, 10, 64); err == nil {
			queryStatement = stmtMemberListByID
		} else if common.RegexpValidateEmailAddress.MatchString(searchQuery) {
			queryStatement = stmtMemberListByEmail
		} else if common.RegexpValidatePostalcode.MatchString(searchQuery) {
			if len(searchQuery) == 7 && searchQuery[4] == ' ' {
				searchQuery = searchQuery[0:3] + searchQuery[5:6]
			}
			searchQuery = strings.ToUpper(searchQuery)
			queryStatement = stmtMemberListByPostalcode
		} else {
			queryStatement = stmtMemberListByTextsearch
			searchQuery = regexpQueryReplacer.ReplaceAllString(searchQuery, "$1 & $2")
		}

		out := &OutMemberList{}
		res, err := queryStatement.Queryx(searchQuery)
		if err != nil {
			if perr, ok := err.(*pq.Error); ok {
				if strings.HasPrefix(perr.Message, `syntax error in tsquery`) ||
					strings.HasPrefix(perr.Message, `no operand in tsquery`) {
					out.Error = `rutte:query_error`
					goto Done
				}
			}
			log.WithError(err).Error("error querying db")
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		for res.Next() {
			var member OutMember
			err := res.StructScan(&member)
			if err != nil {
				log.WithError(err).Error("error scanning member")
				http.Error(w, "server error", http.StatusInternalServerError)
				return
			}
			out.Members = append(out.Members, member)
		}
		if res.Err() != nil {
			log.WithError(res.Err()).Error("error in query result")
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

	Done:
		render.JSON(w, r, out)
	}
}
