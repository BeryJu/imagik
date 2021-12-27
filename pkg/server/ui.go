package server

import (
	"fmt"
	"net/http"

	"beryju.org/imagik/pkg/config"
	"beryju.org/imagik/web/dist"
	"github.com/getsentry/sentry-go"
)

func (s *Server) configureUI() {
	ui := s.handler.NewRoute().Subrouter()
	ui.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			hub := sentry.GetHubFromContext(r.Context())
			hub.Scope().SetTransaction(fmt.Sprintf("%s UI", r.Method))
			h.ServeHTTP(rw, r)
		})
	})
	if config.C.Debug {
		ui.PathPrefix("/ui").Handler(http.StripPrefix("/ui", http.FileServer(http.Dir("./web/dist"))))
	} else {
		ui.PathPrefix("/ui").Handler(http.StripPrefix("/ui", http.FileServer(http.FS(dist.Static))))
	}
	ui.Path("/").HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		http.Redirect(rw, r, "/ui/", http.StatusFound)
	})
}
