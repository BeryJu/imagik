package server

import (
	"net/http"
	"strings"

	"github.com/BeryJu/imagik/root"
)

func (s *Server) UIRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/ui/", http.StatusFound)
}

func (s *Server) UIHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		relativePath := strings.Replace(r.URL.Path, "/ui/", "", 1)
		if relativePath == "" {
			relativePath = "index.html"
		}
		http.FileServer(http.FS(root.Static)).ServeHTTP(w, r)
		if strings.HasSuffix(relativePath, ".js") {
			w.Header().Set("Content-Type", "text/javascript")
		} else if strings.HasSuffix(relativePath, ".css") {
			w.Header().Set("Content-Type", "text/css")
		}
	}
}
