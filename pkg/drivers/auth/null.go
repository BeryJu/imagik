package auth

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

type NullAuth struct {
	logger *log.Entry
}

func (na *NullAuth) Init() {
	na.logger = log.WithField("component", "null-auth")
}

func (na *NullAuth) AuthenticateRequest(w http.ResponseWriter, r *http.Request, next http.Handler) {
	na.logger.Info("Permission denied")
	w.WriteHeader(401)
}
