package server

import (
	"net/http"
	"time"

	"github.com/BeryJu/gopyazo/pkg/config"
	log "github.com/sirupsen/logrus"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		before := time.Now()
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
		after := time.Now()
		log.WithFields(log.Fields{
			"remote": r.RemoteAddr,
			"method": r.Method,
			"took":   after.Sub(before),
		}).Info(r.RequestURI)
	})
}

func configAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config.Config.GetAuth().AuthenticateRequest(w, r, next)
	})
}
