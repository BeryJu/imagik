package auth

import (
	"fmt"
	"net/http"
	"os"

	"github.com/BeryJu/gopyazo/pkg/config"
	"github.com/BeryJu/gopyazo/pkg/drivers"
	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type AuthDriver interface {
	drivers.HTTPDriver
	AuthenticateRequest(w http.ResponseWriter, r *http.Request, next http.Handler)
}

func FromConfig(store *sessions.CookieStore, r *mux.Router) func(next http.Handler) http.Handler {
	authDriverType := config.C.AuthDriver
	var authDriver AuthDriver
	switch authDriverType {
	case "null":
		authDriver = &NullAuth{}
	case "static":
		authDriver = &StaticAuth{}
	case "oidc":
		authDriver = &OIDCAuth{Store: store}
	}
	if authDriver == nil {
		fmt.Printf("Could not configure AuthDriver '%s'", authDriverType)
		os.Exit(1)
	}
	authDriver.Init()
	authDriver.InitRoutes(r)
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			span := sentry.StartSpan(r.Context(), "request.authHandler")
			authDriver.AuthenticateRequest(w, r, next)
			span.Finish()
		})
	}
}
