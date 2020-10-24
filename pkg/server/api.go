package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/BeryJu/gopyazo/pkg/config"
	"github.com/BeryJu/gopyazo/pkg/schema"
	"github.com/vimeo/go-magic/magic"
)

func getElementsForDirectory(path string) int {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return 0
	}
	return len(files)
}

func (s *Server) APIListHandler(w http.ResponseWriter, r *http.Request) {
	offset := r.URL.Query().Get("pathOffset")
	fullDir := config.CleanURL(offset)
	files, err := ioutil.ReadDir(fullDir)
	if err != nil {
		schema.ErrorHandlerAPI(err, w)
		return
	}
	response := schema.ListResponse{
		GenericResponse: schema.GenericResponse{
			Successful: true,
		},
		Files: make([]schema.ListFile, 0),
	}
	for _, f := range files {
		fullName := path.Join(fullDir, f.Name())
		file := schema.ListFile{
			Name:     f.Name(),
			FullPath: fullName,
		}
		if f.IsDir() {
			file.Type = "directory"
			file.ChildElements = getElementsForDirectory(fullName)
		} else {
			file.Type = "file"
			file.Mime = magic.MimeFromFile(fullName)
		}
		response.Files = append(response.Files, file)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		schema.ErrorHandlerAPI(err, w)
		return
	}
}

func (s *Server) APIMoveHandler(w http.ResponseWriter, r *http.Request) {
	fromFull := config.CleanURL(r.URL.Query().Get("from"))
	toFull := config.CleanURL(r.URL.Query().Get("to"))
	if _, err := os.Stat(fromFull); err != nil {
		schema.ErrorHandlerAPI(err, w)
		return
	}
	err := os.Rename(fromFull, toFull)
	if err != nil {
		schema.ErrorHandlerAPI(err, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&schema.GenericResponse{
		Successful: true,
	})
}
