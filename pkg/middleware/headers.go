package middleware

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/csrf"
)

func AddBaseHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-CSRF-TOKEN", csrf.Token(r))
		// AddCorsHeaders(w, r)
		next.ServeHTTP(w, r)
	})
}

func AddJsonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// AddCorsHeaders(w, r)
		next.ServeHTTP(w, r)
	})
}

// func AddCorsHeaders(w http.ResponseWriter, r *http.Request) {
// w.Header().Set("Content-Type", "application/json")
// w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
// w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
// w.Header().Set("Access-Control-Allow-Credentials", "true")
// w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, Cookie, Set-Cookie")
// }

func Csrf() func(http.Handler) http.Handler {
	return csrf.Protect([]byte(os.Getenv("CSRF_AUTH_KEY")), csrf.Path("/"), csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(struct {
			Error string `json:"error"`
		}{"CSRF error"})
	})))

	// return csrfHandler(
	// http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("X-CSRF-Token", csrf.Token(r))
	// next.ServeHTTP(w, r)
	// }),
	// )
}
