package transform

import (
	"encoding/json"
	"net/http"
	"os"

	"beryju.io/imagik/pkg/config"
	"beryju.io/imagik/pkg/drivers/storage"
	"beryju.io/imagik/pkg/schema"
)

type MetaTransformer struct {
	TransformerManager
}

func (mt *MetaTransformer) Handle(w http.ResponseWriter, r *http.Request) {
	fullPath := config.CleanURL(r.URL.Path)
	// Get stat for common file stats
	stat, err := os.Stat(fullPath)
	if err != nil {
		schema.ErrorHandlerAPI(err, w)
		return
	}
	// Get hashes for linking
	hashes, err := mt.sd.HashesForFile(fullPath, storage.ObjectInfo{}, r.Context())
	if err != nil {
		schema.ErrorHandlerAPI(err, w)
		return
	}
	response := schema.MetaResponse{
		GenericResponse: schema.GenericResponse{
			Successful: true,
		},
		Name:         stat.Name(),
		CreationDate: stat.ModTime(),
		Size:         stat.Size(),
		Mime:         hashes.Mime,
		Hashes:       hashes.Map(),
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		schema.ErrorHandlerAPI(err, w)
		return
	}
}
