package schema

type FormUploadResponse struct {
	GenericResponse
	FileResults map[string]string `json:"fileResults"`
}
