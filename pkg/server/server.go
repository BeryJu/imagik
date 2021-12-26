package server

import (
	"fmt"
	"net/http"

	"beryju.org/imagik/pkg/config"
	"beryju.org/imagik/pkg/drivers/auth"
	"beryju.org/imagik/pkg/drivers/metrics"
	"beryju.org/imagik/pkg/hash"
	"beryju.org/imagik/pkg/transform"
	"beryju.org/imagik/root"
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
	md       metrics.MetricsDriver
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
	mainHandler.Use(handlers.ProxyHeaders)
	mainHandler.Use(handlers.CompressHandler)
	mainHandler.Use(loggingMiddleware)
	mainHandler.Use(sentryhttp.New(sentryhttp.Options{}).Handle)
	mainHandler.Use(func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Set("Server", "github.com/beryju/imagik")
			inner.ServeHTTP(rw, r)
		})
	})

	apiPubHandler := mainHandler.PathPrefix("/api/pub").Subrouter()
	authHandler := mainHandler.NewRoute().Subrouter()
	authHandler.Use(auth.FromConfig(store, apiPubHandler))
	apiPrivHandler := authHandler.PathPrefix("/api/priv").Subrouter()
	// apiPrivHandler.Use(csrfMiddleware(apiPrivHandler))

	server.md = metrics.FromConfig(authHandler)

	mainHandler.PathPrefix("/ui").Handler(http.StripPrefix("/ui", http.FileServer(http.FS(root.Static))))
	if !config.C.Debug {
		mainHandler.Path("/").HandlerFunc(server.UIRedirect)
	}
	// General Get Requests don't need authentication
	mainHandler.PathPrefix("/").Methods(http.MethodGet).HandlerFunc(server.GetHandler)
	// Only enable logging middleware after we've added general serving
	authHandler.PathPrefix("/").Methods(http.MethodPut).HandlerFunc(server.PutHandler)
	apiPrivHandler.Path("/list").Methods(http.MethodGet).HandlerFunc(server.APIListHandler)
	apiPrivHandler.Path("/move").Methods(http.MethodPost).HandlerFunc(server.APIMoveHandler)
	apiPrivHandler.Path("/upload").Methods(http.MethodPost).HandlerFunc(server.UploadFormHandler)
	apiPubHandler.Path("/health/liveness").Methods(http.MethodGet).HandlerFunc(server.HealthLiveness)
	apiPubHandler.Path("/health/readiness").Methods(http.MethodGet).HandlerFunc(server.HealthReadiness)

	err := mainHandler.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			server.logger.Debugf("Registered route '%s'", pathTemplate)
		}
		return nil
	})
	if err != nil {
		server.logger.WithError(err).Warning("failed to walk storage")
	}
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
	return http.ListenAndServe(config.C.Listen, s.handler)
}
