package auth

import (
	"net/http"

	"github.com/BeryJu/gopyazo/pkg/drivers"
)

type AuthDriver interface {
	drivers.Driver
	AuthenticateRequest(w http.ResponseWriter, r *http.Request, next http.Handler)
}
