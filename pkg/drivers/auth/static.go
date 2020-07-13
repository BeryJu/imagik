package auth

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type StaticAuth struct {
	staticTokens map[string]string
	logger       *log.Entry
}

func (sa *StaticAuth) Init() {
	sa.staticTokens = viper.GetStringMapString("auth_static_tokens")
	sa.logger = log.WithField("component", "static-auth")
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
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
