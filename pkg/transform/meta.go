package transform

import (
	"encoding/json"
	"net/http"
	"os"

	"beryju.org/imagik/pkg/config"
	"beryju.org/imagik/pkg/hash"
	"beryju.org/imagik/pkg/schema"
	"github.com/gabriel-vasile/mimetype"
)

type MetaTransformer struct {
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
	hashes, err := hash.HashesForFile(fullPath)
	if err != nil {
		schema.ErrorHandlerAPI(err, w)
		return
	}
	mime, err := mimetype.DetectFile(fullPath)
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
		Mime:         mime.String(),
		Hashes:       hashes.Map(),
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		schema.ErrorHandlerAPI(err, w)
		return
	}
}
