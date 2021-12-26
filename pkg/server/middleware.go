package server

import (
	"net/http"
	"time"

	"beryju.org/imagik/pkg/config"
	"github.com/getsentry/sentry-go"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span := sentry.StartSpan(r.Context(), "request.logging")
		// w.Header().Set("Access-Control-Allow-Origin", "*")
		before := time.Now()
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
		after := time.Now()
		log.WithFields(log.Fields{
			"remote": r.RemoteAddr,
			"method": r.Method,
			"took":   after.Sub(before),
		}).Info(r.RequestURI)
		span.Finish()
	})
}

func csrfMiddleware(r *mux.Router) func(next http.Handler) http.Handler {
	csrfMiddleware := csrf.Protect(config.C.SecretKey, csrf.Secure(false))
	r.Use(csrfMiddleware)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			span := sentry.StartSpan(r.Context(), "request.csrf")
			w.Header().Set("X-CSRF-Token", csrf.Token(r))
			next.ServeHTTP(w, r)
			span.Finish()
		})
	}
}
