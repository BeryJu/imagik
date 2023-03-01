package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"beryju.io/imagik/pkg/config"
	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const StaticAuthUser = "imagik_static_user"

type StaticAuth struct {
	Store *sessions.CookieStore

	staticTokens map[string]string
	logger       *log.Entry
}

func (sa *StaticAuth) Init() {
	sa.staticTokens = make(map[string]string, len(config.C.AuthStaticConfig.Tokens))
	for user, pass := range config.C.AuthStaticConfig.Tokens {
		sa.staticTokens[user] = strings.ReplaceAll(pass, "|", "$")
	}

	sa.logger = log.WithField("component", "imagik.drivers.auth.static")
}

func (sa *StaticAuth) InitRoutes(r *mux.Router) {
	r.Path("/auth/driver").HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(rw).Encode(AuthType{
			Type: "static",
		})
	})
	r.Path("/auth/is_authenticated").Methods(http.MethodGet).HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		session, _ := sa.Store.Get(r, SessionName)
		_, ok := session.Values[StaticAuthUser]
		_ = json.NewEncoder(rw).Encode(IsLoggedInResponse{
			Successful: ok,
		})
	})
	r.Path("/auth/login").Methods(http.MethodPost).HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		session, _ := sa.Store.Get(r, SessionName)
		if username, password, found := r.BasicAuth(); found {
			if expectedHash, found := sa.staticTokens[username]; found {
				err := bcrypt.CompareHashAndPassword([]byte(expectedHash), []byte(password))
				if err == nil {
					session.Values[StaticAuthUser] = username
					_ = sa.Store.Save(r, rw, session)
					_ = json.NewEncoder(rw).Encode(IsLoggedInResponse{
						Successful: true,
					})
					return
				}
			}
		}
		_ = json.NewEncoder(rw).Encode(IsLoggedInResponse{
			Successful: false,
		})
	})
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hash := string(bytes)
	return strings.ReplaceAll(hash, "$", "|"), err
}

func (sa *StaticAuth) AuthenticateRequest(w http.ResponseWriter, r *http.Request, next http.Handler) {
	session, _ := sa.Store.Get(r, SessionName)
	if username, ok := session.Values[StaticAuthUser]; ok {
		sa.logger.WithField("user", username).Info("Authenticated as user")
		hub := sentry.GetHubFromContext(r.Context())
		hub.Scope().SetUser(sentry.User{
			Username: username.(string),
		})
		next.ServeHTTP(w, r)
		return
	}
	sa.logger.Info("Permission denied")
	w.WriteHeader(401)
}
