package server

import (
	"net/http"
	"time"

	"github.com/BeryJu/gopyazo/pkg/drivers/auth"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
	authDriverType := viper.GetString("authentication_driver")
	var authDriver auth.AuthDriver
	switch authDriverType {
	case "static":
		authDriver = &auth.StaticAuth{}
	case "null":
		authDriver = &auth.NullAuth{}
	}
	authDriver.Init()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authDriver.AuthenticateRequest(w, r, next)
	})
}
