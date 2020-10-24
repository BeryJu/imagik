package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/BeryJu/gopyazo/pkg/schema"
	"github.com/vimeo/go-magic/magic"
)

func errorHandlerAPI(err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schema.GenericResponse{
		Successful: false,
		Error:      err.Error(),
	})
}

func (s *Server) APIListHandler(w http.ResponseWriter, r *http.Request) {
	offset := r.URL.Query().Get("pathOffset")
	fullDir := s.cleanURL(offset)
	files, err := ioutil.ReadDir(fullDir)
	if err != nil {
		errorHandlerAPI(err, w)
		return
	}
	response := schema.ListResponse{
		GenericResponse: schema.GenericResponse{
			Successful: true,
		},
		Directories: make([]schema.ListDirectory, 0),
		Files:       make([]schema.ListFile, 0),
	}
	for _, f := range files {
		if f.IsDir() {
			dir := schema.ListDirectory{Name: f.Name()}
			response.Directories = append(response.Directories, dir)
		} else {
			file := schema.ListFile{
				Name: f.Name(),
				Mime: magic.MimeFromFile(path.Join(fullDir, f.Name())),
			}
			response.Files = append(response.Files, file)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		errorHandlerAPI(err, w)
		return
	}
}

func (s *Server) APIMoveHandler(w http.ResponseWriter, r *http.Request) {
	fromFull := s.cleanURL(r.URL.Query().Get("from"))
	toFull := s.cleanURL(r.URL.Query().Get("to"))
	if _, err := os.Stat(fromFull); err != nil {
		errorHandlerAPI(err, w)
		return
	}
	err := os.Rename(fromFull, toFull)
	if err != nil {
		errorHandlerAPI(err, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&schema.GenericResponse{
		Successful: true,
	})
}
