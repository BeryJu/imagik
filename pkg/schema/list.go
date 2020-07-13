package schema

type ListResponse struct {
	GenericResponse
	Directories []ListDirectory `json:"directories"`
	Files       []ListFile      `json:"files"`
}

type ListFile struct {
	Name string `json:"name"`
}
type ListDirectory struct {
	Name string `json:"name"`
}
