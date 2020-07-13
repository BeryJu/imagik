package auth

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type StaticAuth struct {
	staticTokens map[string]string
	logger       *log.Entry
}

func (sa *StaticAuth) Init(driverConfig map[string]string) {
	sa.staticTokens = driverConfig
	sa.logger = log.WithField("component", "static-auth")
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
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
	sa.logger.Info("Permission denied")
	w.WriteHeader(401)
}
