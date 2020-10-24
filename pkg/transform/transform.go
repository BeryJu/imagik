package transform

import "net/http"

type Transformer interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type TransformerManager struct {
	transformers map[string]Transformer
}

func New() *TransformerManager {
	return &TransformerManager{
		transformers: map[string]Transformer{
			"meta": &MetaTransformer{},
		},
	}
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
