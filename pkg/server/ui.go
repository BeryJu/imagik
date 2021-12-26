package server

import (
	"net/http"
)

func (s *Server) UIRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/ui/", http.StatusFound)
}
