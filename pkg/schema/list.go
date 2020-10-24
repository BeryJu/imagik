package schema

type ListResponse struct {
	GenericResponse
	Files []ListFile `json:"files"`
}

type ListFile struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	Mime          string `json:"mime"`
	FullPath      string `json:"fullPath"`
	ChildElements int    `json:"childElements"`
}
