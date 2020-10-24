package schema

type ListResponse struct {
	GenericResponse
	Directories []ListDirectory `json:"directories"`
	Files       []ListFile      `json:"files"`
}

type ListFile struct {
	Name string `json:"name"`
	Mime string `json:"mime"`
}

type ListDirectory struct {
	Name          string `json:"name"`
	ChildElements int    `json:"childElements"`
}
