package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/mux"
	"github.com/michaeltintiuc/shackle-api/pkg/utils"
)

func Auth(jwtSecret string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := getToken(r)

			if err != nil {
				utils.HasError(w, err, "JWT token not found", http.StatusForbidden)
				return
			}

			parsedToken, err := jwt.ParseWithClaims(
				token,
				&utils.JWTClaims{},
				func(token *jwt.Token) (interface{}, error) {
					return []byte(jwtSecret), nil
				},
			)

			if err != nil {
				utils.HasError(w, err, "Invalid JWT token", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), "claims", parsedToken.Claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func getToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil
	}

	authHeaderParts := strings.Fields(authHeader)
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", errors.New("Authorization header format must be Bearer {token}")
	}

	return authHeaderParts[1], nil
}
