package schema

import "time"

type MetaResponse struct {
	GenericResponse
	CreationDate time.Time         `json:"creationData"`
	Size         int64             `json:"size"`
	Hashes       map[string]string `json:"hashes"`
	Mime         string            `json:"mime"`
}
