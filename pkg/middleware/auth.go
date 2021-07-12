package middleware

import (
	"encoding/json"
	"log"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/mux"
)

func Auth(jwtSecret string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		j := jwtmiddleware.New(jwtmiddleware.Options{
			ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			},
			ErrorHandler: func(w http.ResponseWriter, r *http.Request, err string) {
				log.Println(err)
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(struct {
					Error string `json:"error"`
				}{"Invalid JWT"})
			},
			SigningMethod: jwt.SigningMethodHS256,
		})

		return j.Handler(next)
	}
}
