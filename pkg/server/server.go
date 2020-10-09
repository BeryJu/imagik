package server

import (
	"fmt"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/BeryJu/gopyazo/pkg/config"
	"github.com/BeryJu/gopyazo/pkg/hash"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	rootDir string
	handler *mux.Router
	logger  *log.Entry
	HashMap *hash.HashMap
}

func New() *Server {
	mainHandler := mux.NewRouter()
	server := &Server{
		rootDir: config.C.RootDir,
		handler: mainHandler,
		logger:  log.WithField("component", "server"),
	}
	mainHandler.Use(handlers.ProxyHeaders)
	mainHandler.Use(loggingMiddleware)
	mainHandler.Use(handlers.CompressHandler)

	apiPubHandler := mainHandler.PathPrefix("/api/pub").Subrouter()
	authHandler := mainHandler.NewRoute().Subrouter()
	authHandler.Use(configAuthMiddleware(apiPubHandler))
	apiPrivHandler := authHandler.PathPrefix("/api/priv").Subrouter()

	// General Get Requests don't need authentication
	mainHandler.PathPrefix("/").Methods(http.MethodGet).HandlerFunc(server.GetHandler)
	authHandler.PathPrefix("/").Methods(http.MethodPut).HandlerFunc(server.PutHandler)
	apiPrivHandler.Path("/list").HandlerFunc(server.APIListHandler)
	apiPrivHandler.Path("/move").HandlerFunc(server.APIMoveHandler)
	apiPubHandler.Path("/health/liveness").HandlerFunc(server.HealthLiveness)
	apiPubHandler.Path("/health/readiness").HandlerFunc(server.HealthReadiness)

	mainHandler.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			server.logger.Debugf("Registered route '%s'", pathTemplate)
		}
		return nil
	})
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

func (s *Server) Run() error {
	log.WithField("listen", config.C.Listen).Info("Server running")
	sentryHandler := sentryhttp.New(sentryhttp.Options{})
	return http.ListenAndServe(config.C.Listen, sentryHandler.Handle(s.handler))
}
