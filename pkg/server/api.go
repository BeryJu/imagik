package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/BeryJu/gopyazo/pkg/hash"
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

func getElementsForDirectory(path string) int {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return 0
	}
	return len(files)
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
		fullName := path.Join(fullDir, f.Name())
		if f.IsDir() {
			dir := schema.ListDirectory{
				Name:          f.Name(),
				ChildElements: getElementsForDirectory(fullName),
			}
			response.Directories = append(response.Directories, dir)
		} else {
			file := schema.ListFile{
				Name: f.Name(),
				Mime: magic.MimeFromFile(fullName),
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

func (s *Server) APIMetaHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	fullPath := s.cleanURL(path)
	// Get stat for common file stats
	stat, err := os.Stat(fullPath)
	if err != nil {
		errorHandlerAPI(err, w)
		return
	}
	// Get hashes for linking
	hashes, err := hash.HashesForFile(fullPath)
	if err != nil {
		errorHandlerAPI(err, w)
		return
	}
	response := schema.MetaResponse{
		GenericResponse: schema.GenericResponse{
			Successful: true,
		},
		CreationDate: stat.ModTime(),
		Size:         stat.Size(),
		Mime:         magic.MimeFromFile(fullPath),
		Hashes:       hashes.Map(),
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
