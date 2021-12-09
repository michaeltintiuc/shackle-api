package middleware

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/michaeltintiuc/shackle-api/pkg/session"
	"github.com/michaeltintiuc/shackle-api/pkg/utils"
)

func Auth(sessionInfo *session.Session) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "OPTIONS" {
				next.ServeHTTP(w, r)
				return
			}

			session, err := sessionInfo.Get(r)
			if utils.HasError(w, err, "Application error", http.StatusInternalServerError) {
				return
			}

			id, found := session.Values["id"]
			if !found || id == "" {
				msg := "Unauthenticated access"
				utils.HasError(w, errors.New(msg), msg, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
