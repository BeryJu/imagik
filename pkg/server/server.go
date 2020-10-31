package server

import (
	"fmt"
	"net/http"

	"github.com/BeryJu/gopyazo/pkg/config"
	"github.com/BeryJu/gopyazo/pkg/hash"
	"github.com/BeryJu/gopyazo/pkg/transform"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	rootDir  string
	handler  *mux.Router
	logger   *log.Entry
	HashMap  *hash.HashMap
	tm       *transform.TransformerManager
	sessions *sessions.CookieStore
}

func New() *Server {
	store := sessions.NewCookieStore(config.C.SecretKey)

	mainHandler := mux.NewRouter()
	server := &Server{
		rootDir:  config.C.RootDir,
		handler:  mainHandler,
		logger:   log.WithField("component", "server"),
		tm:       transform.New(),
		sessions: store,
	}
	mainHandler.Use(recoveryMiddleware())
	mainHandler.Use(handlers.ProxyHeaders)
	mainHandler.Use(loggingMiddleware)
	mainHandler.Use(handlers.CompressHandler)

	apiPubHandler := mainHandler.PathPrefix("/api/pub").Subrouter()
	authHandler := mainHandler.NewRoute().Subrouter()
	authHandler.Use(configAuthMiddleware(store, apiPubHandler))
	apiPrivHandler := authHandler.PathPrefix("/api/priv").Subrouter()
	apiPrivHandler.Use(csrfMiddleware(apiPrivHandler))

	mainHandler.PathPrefix("/ui").HandlerFunc(server.UIHandler())
	if !config.C.Debug {
		mainHandler.Path("/").HandlerFunc(server.UIRedirect)
	}
	// General Get Requests don't need authentication
	mainHandler.PathPrefix("/").Methods(http.MethodGet).HandlerFunc(server.GetHandler)
	authHandler.PathPrefix("/").Methods(http.MethodPut).HandlerFunc(server.PutHandler)
	apiPrivHandler.Path("/list").Methods(http.MethodGet).HandlerFunc(server.APIListHandler)
	apiPrivHandler.Path("/move").Methods(http.MethodPost).HandlerFunc(server.APIMoveHandler)
	apiPrivHandler.Path("/upload").Methods(http.MethodPost).HandlerFunc(server.UploadFormHandler)
	apiPubHandler.Path("/health/liveness").Methods(http.MethodGet).HandlerFunc(server.HealthLiveness)
	apiPubHandler.Path("/health/readiness").Methods(http.MethodGet).HandlerFunc(server.HealthReadiness)

	mainHandler.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			server.logger.Debugf("Registered route '%s'", pathTemplate)
		}
		return nil
	})
	return server
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
