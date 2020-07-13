package server

import (
	"io"
	"net/http"
	"os"
	"path"
)

// GetHandler Handle GET Requests
func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	filePath := s.cleanURL(r.URL.Path)
	_, err := os.Stat(filePath)
	if err == nil {
		s.logger.Debug("Handling normal serve")
		http.ServeFile(w, r, filePath)
		return
	}
	// Since we only store the hash, we need to get rid of the lading slash
	path, exists := s.HashMap.Get(r.URL.Path[1:])
	if exists {
		s.logger.WithField("path", path).Debug("Found path in hashmap")
		http.ServeFile(w, r, path)
		return
	}
	s.logger.Debug("Not found in hashmap")
	notFoundHandler("File not found.", w)
}

func (s *Server) PutHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	filePath := s.cleanURL(r.URL.Path)
	err := os.MkdirAll(path.Dir(filePath), os.ModePerm)
	if err != nil {
		s.logger.Warning(err)
		errorHandler(err, w)
		return
	}

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		s.logger.Warning(err)
		errorHandler(err, w)
		return
	}
	n, err := io.Copy(f, r.Body)
	if err != nil {
		s.logger.Warning(err)
		errorHandler(err, w)
		return
	}
	s.logger.WithField("n", n).WithField("path", filePath).Debug("Successfully wrote file.")
	s.HashMap.UpdateSingle(filePath)
}
