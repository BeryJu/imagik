package auth

import (
	"net/http"
	"strings"

	"github.com/BeryJu/gopyazo/pkg/config"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type StaticAuth struct {
	staticTokens map[string]string
	logger       *log.Entry
}

func (sa *StaticAuth) Init() {
	sa.staticTokens = make(map[string]string, len(config.C.AuthStaticConfig.Tokens))
	for user, pass := range config.C.AuthStaticConfig.Tokens {
		sa.staticTokens[user] = strings.ReplaceAll(pass, "|", "$")
	}

	sa.logger = log.WithField("component", "static-auth")
}

func (sa *StaticAuth) InitRoutes(r *mux.Router) {
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hash := string(bytes)
	return strings.ReplaceAll(hash, "$", "|"), err
}

func (sa *StaticAuth) AuthenticateRequest(w http.ResponseWriter, r *http.Request, next http.Handler) {
	if username, password, found := r.BasicAuth(); found {
		if expectedHash, found := sa.staticTokens[username]; found {
			err := bcrypt.CompareHashAndPassword([]byte(expectedHash), []byte(password))
			if err == nil {
				sa.logger.WithField("user", username).Info("Authenticated as user")
				next.ServeHTTP(w, r)
				return
			}
		}
	}
	w.Header().Set("WWW-Authenticate", `Basic realm="gopyazo"`)
	sa.logger.Info("Permission denied")
	w.WriteHeader(401)
}
