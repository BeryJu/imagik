package server

import (
	"net/http"

	"beryju.org/imagik/pkg/config"
	"beryju.org/imagik/web/dist"
)

func (s *Server) configureUI() {
	if config.C.Debug {
		s.handler.PathPrefix("/ui").Handler(http.StripPrefix("/ui", http.FileServer(http.Dir("./web/dist"))))
	} else {
		s.handler.PathPrefix("/ui").Handler(http.StripPrefix("/ui", http.FileServer(http.FS(dist.Static))))
	}
	s.handler.Path("/").HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		http.Redirect(rw, r, "/ui/", http.StatusFound)
	})
}
