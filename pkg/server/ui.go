package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gobuffalo/packr/v2"
)

func (s *Server) UIRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/ui/", http.StatusFound)
}

func (s *Server) UIHandler() func(w http.ResponseWriter, r *http.Request) {
	uiBox := packr.New("webui", "../../root")
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
		fmt.Println(uiBox.List())
		relativePath := strings.Replace(r.URL.Path, "/ui/", "", 1)
		fmt.Println(relativePath)
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
