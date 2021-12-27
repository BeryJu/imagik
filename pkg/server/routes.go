package server

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"beryju.org/imagik/pkg/config"
	"beryju.org/imagik/pkg/drivers/metrics"
	"beryju.org/imagik/pkg/hash"
	"beryju.org/imagik/pkg/schema"
	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
)

// GetHandler Handle GET Requests
func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	span := sentry.StartSpan(r.Context(), "request.GetFile")
	defer span.Finish()

	if s.tm.Transform(w, r) {
		return
	}
	filePath := config.CleanURL(r.URL.Path)
	mr := metrics.NewServeRequest(r)
	start := time.Now()
	defer func() {
		mr.Duration = time.Since(start)
	}()
	// Since we only store the hash, we need to get rid of the leading slash
	p, exists := s.HashMap.Get(r.URL.Path[1:])
	if exists {
		mr.Hash = r.URL.Path[1:]
		mr.ResolvedPath = p
		s.md.ServeRequest(mr)
		http.ServeFile(w, r, p)
		return
	}
	// Check if we have the file without extension
	base := path.Base(r.URL.Path[1:])
	ext := path.Ext(base)
	filename := strings.Replace(base, ext, "", 1)
	p, exists = s.HashMap.Get(filename)
	if exists {
		mr.Hash = r.URL.Path[1:]
		mr.ResolvedPath = p
		s.md.ServeRequest(mr)
		http.ServeFile(w, r, p)
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
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		s.logger.WithError(err).Warning("failed to parse multipart form")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
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
			_, err := s.doUpload(mf, key)
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
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		s.logger.WithError(err).Warning("failed to write json response")
	}
}

// PutHandler Upload handler used frm CLI
func (s *Server) PutHandler(w http.ResponseWriter, r *http.Request) {
	span := sentry.StartSpan(r.Context(), "request.PutFile")
	defer span.Finish()
	hashes, err := s.doUpload(r.Body, r.URL.Path)
	if err != nil {
		errorHandler(err, w)
		return
	}
	w.WriteHeader(200)
	err = json.NewEncoder(w).Encode(hashes)
	if err != nil {
		schema.ErrorHandlerAPI(err, w)
		return
	}
}

func (s *Server) doUpload(src io.Reader, p string) (*hash.FileHash, error) {
	filePath := config.CleanURL(p)
	err := os.MkdirAll(path.Dir(filePath), os.ModePerm)
	if err != nil {
		s.logger.Warning(err)
		return nil, err
	}

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		s.logger.Warning(err)
		return nil, err
	}
	n, err := io.Copy(f, src)
	if err != nil {
		s.logger.Warning(err)
		return nil, err
	}
	s.logger.WithField("n", n).WithField("path", filePath).Debug("Successfully wrote file.")
	return s.HashMap.UpdateSingle(filePath), nil
}
