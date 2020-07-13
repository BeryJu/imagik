package schema

type GenericResponse struct {
	Successful bool   `json:"successful"`
	Error      string `json:"error"`
}
