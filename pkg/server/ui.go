package server

import (
	"net/http"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gobuffalo/packr/v2"
)

func (s *Server) UIHandler() func(w http.ResponseWriter, r *http.Request) {
	uiBox := packr.New("webui", "../../root")
	return func(w http.ResponseWriter, r *http.Request) {
		relativePath := strings.Replace(r.URL.Path, "/ui/", "", 1)
		if relativePath == "" {
			relativePath = "index.html"
		}
		body, err := uiBox.Find(relativePath)
		if strings.HasSuffix(relativePath, ".js") {
			w.Header().Set("Content-Type", "text/javascript")
		} else if strings.HasSuffix(relativePath, ".css") {
			w.Header().Set("Content-Type", "text/css")
		} else {
			w.Header().Set("Content-Type", mimetype.Detect(body).String())
		}
		if err != nil {
			errorHandler(err, w)
			return
		}
		w.WriteHeader(200)
		w.Write(body)
	}
}
