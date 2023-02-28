package transform

import (
	"encoding/json"
	"net/http"

	"beryju.io/imagik/pkg/schema"
)

type MetaTransformer struct {
	TransformerManager
}

func (mt *MetaTransformer) Handle(w http.ResponseWriter, r *http.Request) {
	fullPath := mt.sd.CleanURL(r.URL.Path)
	response := mt.sd.Stat(fullPath, r.Context())
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		schema.ErrorHandlerAPI(err, w)
		return
	}
}
