package server

import (
	"fmt"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/BeryJu/gopyazo/pkg/config"
	"github.com/BeryJu/gopyazo/pkg/hash"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	handler *mux.Router
	logger  *log.Entry
	HashMap *hash.HashMap
}

func New() *Server {
	handler := mux.NewRouter()
	server := &Server{
		handler: handler,
		logger:  log.WithField("component", "server"),
	}
	handler.Use(loggingMiddleware)
	authenticationSubRouter := handler.NewRoute().Subrouter()
	authenticationSubRouter.Use(configAuthMiddleware)
	handler.PathPrefix("/").Methods(http.MethodGet).HandlerFunc(server.GetHandler)
	authenticationSubRouter.PathPrefix("/").Methods(http.MethodPut).HandlerFunc(server.PutHandler)
	return server
}

func (s *Server) cleanURL(raw string) string {
	if !strings.HasPrefix(raw, "/") {
		raw = "/" + raw
	}
	return filepath.Join(config.Config.RootDir, filepath.FromSlash(path.Clean("/"+raw)))
}

func errorHandler(err error, w http.ResponseWriter) {
	fmt.Fprintf(w, "Error: %s", err)
}

func notFoundHandler(msg string, w http.ResponseWriter) {
	w.WriteHeader(404)
	fmt.Fprint(w, msg)
}

func (s *Server) Run() {
	listen := "localhost:8000"
	log.Infof("Server running '%s'", listen)
	http.ListenAndServe(listen, s.handler)
}
