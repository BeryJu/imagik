package ui

import (
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"beryju.io/imagik/pkg/config"
	"beryju.io/imagik/web/dist"
	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
)

//go:embed index.html
var IndexTemplate string

func ConfigureUI(h *mux.Router) {
	ui := h.NewRoute().Subrouter()
	ui.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			hub := sentry.GetHubFromContext(r.Context())
			hub.Scope().SetTransaction(fmt.Sprintf("%s UI", r.Method))
			h.ServeHTTP(rw, r)
		})
	})
	t, err := template.New("imagik.ui").Parse(IndexTemplate)
	if err != nil {
		log.Fatalf("failed parsing template %s", err)
	}

	ui.Path("/ui/").HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		span := sentry.StartSpan(r.Context(), "imagik.ui")
		defer span.Finish()
		er := t.Execute(rw, struct {
			Trace string
		}{
			Trace: span.ToSentryTrace(),
		})
		if er != nil {
			http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		}
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
