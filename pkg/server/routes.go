package server

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/BeryJu/imagik/pkg/config"
	"github.com/BeryJu/imagik/pkg/drivers/metrics"
	"github.com/BeryJu/imagik/pkg/schema"
	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
)

// GetHandler Handle GET Requests
func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	span := sentry.StartSpan(r.Context(), "request.GetFile")
	defer span.Finish()

	if s.tm.Transform(w, r) == true {
		return
	}
	filePath := config.CleanURL(r.URL.Path)
	mr := metrics.NewServeRequest(r)
	// Since we only store the hash, we need to get rid of the leading slash
	path, exists := s.HashMap.Get(r.URL.Path[1:])
	if exists {
		mr.Hash = r.URL.Path[1:]
		mr.ResolvedPath = path
		s.md.ServeRequest(mr)
		http.ServeFile(w, r, path)
		return
	}
	_, err := os.Stat(filePath)
	if err == nil {
		mr.ResolvedPath = filePath
		s.md.ServeRequest(mr)
		http.ServeFile(w, r, filePath)
		return
	}
	notFoundHandler("File not found.", w)
}

// UploadFormHandler Upload handler used by HTML Forms
func (s *Server) UploadFormHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	fileResultMap := make(map[string]string, len(r.MultipartForm.File))
	for key, files := range r.MultipartForm.File {
		if len(files) < 1 {
			continue
		}
		file := files[0]
		mf, err := file.Open()
		if err != nil {
			fileResultMap[key] = errors.Wrap(err, "failed to open multipart file").Error()
		} else {
			err := s.doUpload(mf, key)
			if err != nil {
				fileResultMap[key] = err.Error()
			} else {
				fileResultMap[key] = ""
			}
		}
	}
	response := schema.FormUploadResponse{
		GenericResponse: schema.GenericResponse{
			Successful: true,
		},
		FileResults: fileResultMap,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// PutHandler Upload handler used frm CLI
func (s *Server) PutHandler(w http.ResponseWriter, r *http.Request) {
	span := sentry.StartSpan(r.Context(), "request.PutFile")
	defer span.Finish()
	err := s.doUpload(r.Body, r.URL.Path)
	if err != nil {
		errorHandler(err, w)
		return
	}
	w.WriteHeader(201)
}

func (s *Server) doUpload(src io.Reader, p string) error {
	filePath := config.CleanURL(p)
	err := os.MkdirAll(path.Dir(filePath), os.ModePerm)
	if err != nil {
		s.logger.Warning(err)
		return err
	}

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		s.logger.Warning(err)
		return err
	}
	n, err := io.Copy(f, src)
	if err != nil {
		s.logger.Warning(err)
		return err
	}
	s.logger.WithField("n", n).WithField("path", filePath).Debug("Successfully wrote file.")
	s.HashMap.UpdateSingle(filePath)
	return nil
}
