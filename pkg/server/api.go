package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"beryju.org/imagik/pkg/config"
	"beryju.org/imagik/pkg/schema"
	"github.com/gabriel-vasile/mimetype"
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
			FullPath: filepath.Join(filepath.FromSlash(path.Clean("/"+offset)), f.Name()),
		}
		if f.IsDir() {
			file.Type = "directory"
			file.ChildElements = getElementsForDirectory(fullName)
		} else {
			file.Type = "file"
			mime, err := mimetype.DetectFile(fullName)
			if err == nil {
				file.Mime = mime.String()
			}
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
	var from, to string
	if from = r.URL.Query().Get("from"); from == "" {
		schema.ErrorHandlerAPI(errors.New("missing from path"), w)
		return
	}
	if to = r.URL.Query().Get("to"); to == "" {
		schema.ErrorHandlerAPI(errors.New("missing to path"), w)
		return
	}
	fromFull := config.CleanURL(from)
	toFull := config.CleanURL(to)
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
	err = json.NewEncoder(w).Encode(&schema.GenericResponse{
		Successful: true,
	})
	if err != nil {
		s.logger.WithError(err).Warning("failed to write json response")
	}
}
