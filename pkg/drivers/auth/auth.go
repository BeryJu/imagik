package auth

import (
	"fmt"
	"net/http"
	"os"

	"beryju.org/imagik/pkg/config"
	"beryju.org/imagik/pkg/drivers"
	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

const SessionName = "imagik_session"

type AuthType struct {
	Type string            `json:"type"`
	Args map[string]string `json:"args"`
}

type IsLoggedInResponse struct {
	Successful bool
}

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
		authDriver = &StaticAuth{Store: store}
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
