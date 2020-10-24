package transform

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/BeryJu/gopyazo/pkg/config"
	"github.com/BeryJu/gopyazo/pkg/hash"
	"github.com/BeryJu/gopyazo/pkg/schema"
	"github.com/vimeo/go-magic/magic"
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
	response := schema.MetaResponse{
		GenericResponse: schema.GenericResponse{
			Successful: true,
		},
		CreationDate: stat.ModTime(),
		Size:         stat.Size(),
		Mime:         magic.MimeFromFile(fullPath),
		Hashes:       hashes.Map(),
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		schema.ErrorHandlerAPI(err, w)
		return
	}
}
