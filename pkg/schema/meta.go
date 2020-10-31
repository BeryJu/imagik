package schema

import "time"

type MetaResponse struct {
	GenericResponse
	Name         string            `json:"name"`
	CreationDate time.Time         `json:"creationDate"`
	Size         int64             `json:"size"`
	Hashes       map[string]string `json:"hashes"`
	Mime         string            `json:"mime"`
}
