package schema

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type GenericResponse struct {
	Successful bool   `json:"successful"`
	Error      string `json:"error"`
}

func ErrorHandlerAPI(err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	err = json.NewEncoder(w).Encode(GenericResponse{
		Successful: false,
		Error:      err.Error(),
	})
	if err != nil {
		log.WithError(err).Warning("failed to write json error")
	}
}
