package middleware

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type statusRecorder struct {
	status int
	w      http.ResponseWriter
}

func (srw *statusRecorder) Header() http.Header {
	return srw.w.Header()
}

func (srw *statusRecorder) Write(p []byte) (int, error) {
	if srw.status == 0 {
		srw.status = http.StatusOK
	}
	return srw.w.Write(p)
}

func (srw *statusRecorder) WriteHeader(status int) {
	srw.status = status
	srw.w.WriteHeader(status)
}

func Init(r *mux.Router) {
	r.Use(logger)
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srw := &statusRecorder{w: w}
		defer func() {
			log.Printf("[%d] %s %s - %s", srw.status, r.Method, r.RemoteAddr, r.RequestURI)
		}()
		next.ServeHTTP(srw, r)
	})
}
