package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"beryju.io/imagik/pkg/schema"
	"github.com/getsentry/sentry-go"
)

func (s *Server) APIListHandler(w http.ResponseWriter, r *http.Request) {
	tx := sentry.TransactionFromContext(r.Context())
	if tx != nil {
		tx.Name = fmt.Sprintf("%s APIList", r.Method)
	}
	offset := r.URL.Query().Get("pathOffset")
	files, err := s.sd.List(r.Context(), offset)
	if err != nil {
		schema.ErrorHandlerAPI(err, w)
		return
	}
	response := schema.ListResponse{
		GenericResponse: schema.GenericResponse{
			Successful: true,
		},
		Files: files,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		schema.ErrorHandlerAPI(err, w)
		return
	}
}

func (s *Server) APIMoveHandler(w http.ResponseWriter, r *http.Request) {
	tx := sentry.TransactionFromContext(r.Context())
	if tx != nil {
		tx.Name = fmt.Sprintf("%s APIMove", r.Method)
	}
	var from, to string
	if from = r.URL.Query().Get("from"); from == "" {
		schema.ErrorHandlerAPI(errors.New("missing from path"), w)
		return
	}
	if to = r.URL.Query().Get("to"); to == "" {
		schema.ErrorHandlerAPI(errors.New("missing to path"), w)
		return
	}
	fromFull := s.sd.CleanURL(from)
	toFull := s.sd.CleanURL(to)
	err := s.sd.Rename(r.Context(), fromFull, toFull)
	if err != nil {
		schema.ErrorHandlerAPI(err, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&schema.GenericResponse{
		Successful: true,
	})
	if err != nil {
		s.logger.WithError(err).Warning("failed to write json response")
	}
}
