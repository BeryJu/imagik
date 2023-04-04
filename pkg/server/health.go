package server

import (
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go"
)

func (s *Server) HealthLiveness(w http.ResponseWriter, r *http.Request) {
	tx := sentry.TransactionFromContext(r.Context())
	if tx != nil {
		tx.Name = fmt.Sprintf("%s HealthLiveness", r.Method)
	}
	w.WriteHeader(201)
}

func (s *Server) HealthReadiness(w http.ResponseWriter, r *http.Request) {
	tx := sentry.TransactionFromContext(r.Context())
	if tx != nil {
		tx.Name = fmt.Sprintf("%s HealthReadiness", r.Method)
	}
	if s.HashMap.Populated() {
		w.WriteHeader(201)
	} else {
		w.WriteHeader(500)
	}
}
