package middleware

import (
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/michaeltintiuc/shackle-api/pkg/utils"
)

func Auth(jwtSecret string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			signedToken, err := jwtmiddleware.FromAuthHeader(r)

			if err != nil {
				utils.HasError(w, err, "JWT token not found", http.StatusForbidden)
				return
			}

			_, err = jwt.ParseWithClaims(
				signedToken,
				&utils.JWTClaims{},
				func(token *jwt.Token) (interface{}, error) {
					return []byte(jwtSecret), nil
				},
			)

			if err != nil {
				utils.HasError(w, err, "Invalid JWT token", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
