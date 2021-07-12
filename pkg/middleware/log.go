package middleware

import (
	"log"
	"net/http"
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

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srw := &statusRecorder{w: w}
		next.ServeHTTP(srw, r)
		log.Printf("[%d] %s %s - %s\n", srw.status, r.Method, r.RemoteAddr, r.RequestURI)
	})
}
