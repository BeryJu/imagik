package transform

import (
	"net/http"

	"beryju.io/imagik/pkg/drivers/storage"
)

type Transformer interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type TransformerManager struct {
	transformers map[string]Transformer
	sd           storage.Driver
}

func New(sd storage.Driver) *TransformerManager {
	tm := &TransformerManager{
		sd:           sd,
		transformers: make(map[string]Transformer),
	}
	tm.transformers["meta"] = &MetaTransformer{*tm}
	return tm
}

func (tm *TransformerManager) Transform(w http.ResponseWriter, r *http.Request) bool {
	for key, transformer := range tm.transformers {
		if _, keySet := r.URL.Query()[key]; keySet {
			transformer.Handle(w, r)
			return true
		}
	}
	return false
}
