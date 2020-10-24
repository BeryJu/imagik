package schema

import (
	"encoding/json"
	"net/http"
)

type GenericResponse struct {
	Successful bool   `json:"successful"`
	Error      string `json:"error"`
}

func ErrorHandlerAPI(err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GenericResponse{
		Successful: false,
		Error:      err.Error(),
	})
}
