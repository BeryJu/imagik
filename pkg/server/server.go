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
	"github.com/spf13/viper"
)

type Server struct {
	rootDir string
	handler *mux.Router
	logger  *log.Entry
	HashMap *hash.HashMap
}

func New() *Server {
	handler := mux.NewRouter()
	server := &Server{
		rootDir: viper.GetString(config.ConfigRootDir),
		handler: handler,
		logger:  log.WithField("component", "server"),
	}
	handler.Use(loggingMiddleware)

	authenticationSubRouter := handler.NewRoute().Subrouter()
	authenticationSubRouter.Use(configAuthMiddleware)
	apiRouter := authenticationSubRouter.PathPrefix(viper.GetString(config.ConfigAPIPathPrefix)).Subrouter()

	// General Get Requests don't need authentication
	handler.PathPrefix("/").Methods(http.MethodGet).HandlerFunc(server.GetHandler)
	authenticationSubRouter.PathPrefix("/").Methods(http.MethodPut).HandlerFunc(server.PutHandler)
	apiRouter.Path("/list").HandlerFunc(server.APIListHandler)
	apiRouter.Path("/move").HandlerFunc(server.APIMoveHandler)
	return server
}

func (s *Server) cleanURL(raw string) string {
	if !strings.HasPrefix(raw, "/") {
		raw = "/" + raw
	}
	return filepath.Join(s.rootDir, filepath.FromSlash(path.Clean("/"+raw)))
}

func errorHandler(err error, w http.ResponseWriter) {
	fmt.Fprintf(w, "Error: %s", err)
}

func notFoundHandler(msg string, w http.ResponseWriter) {
	w.WriteHeader(404)
	fmt.Fprint(w, msg)
}

func (s *Server) Run() {
	log.Infof("Server running '%s'", viper.GetString(config.ConfigListen))
	http.ListenAndServe(viper.GetString(config.ConfigListen), s.handler)
}
