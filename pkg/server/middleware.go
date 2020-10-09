package server

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/BeryJu/gopyazo/pkg/config"
	"github.com/BeryJu/gopyazo/pkg/drivers/auth"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	})
}

func configAuthMiddleware(r *mux.Router) func(next http.Handler) http.Handler {
	authDriverType := viper.GetString(config.ConfigAuthenticationDriver)
	var authDriver auth.AuthDriver
	switch authDriverType {
	case "static":
		authDriver = &auth.StaticAuth{}
	case "null":
		authDriver = &auth.NullAuth{}
	case "oidc":
		authDriver = &auth.OIDCAuth{}
	}
	if authDriver == nil {
		fmt.Printf("Could not configure AuthDriver '%s'", authDriverType)
		os.Exit(1)
	}
	authDriver.Init()
	authDriver.InitRoutes(r)
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authDriver.AuthenticateRequest(w, r, next)
		})
	}
}
