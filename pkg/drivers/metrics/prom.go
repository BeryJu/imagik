package metrics

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	Requests = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "imagik_requests",
		Help: "Total requests",
	}, []string{"host", "path", "hash", "client"})
)

type PrometheusMetricsDriver struct {
	m      *mux.Router
	logger *log.Entry
}

func (pmd *PrometheusMetricsDriver) Init() {
	pmd.m = mux.NewRouter()
	pmd.m.Path("/metrics").HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		promhttp.InstrumentMetricHandler(
			prometheus.DefaultRegisterer, promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{
				DisableCompression: true,
			}),
		).ServeHTTP(rw, r)
	})
	pmd.logger = log.WithField("component", "imagik.drivers.metrics.prometheus")
	go pmd.start()
}

func (pmd *PrometheusMetricsDriver) start() {
	listen := "0.0.0.0:9300"
	pmd.logger.WithField("listen", listen).Info("Starting Metrics server")
	err := http.ListenAndServe(listen, pmd.m)
	if err != nil {
		pmd.logger.WithError(err).Warning("Failed to start metrics server")
	}
	pmd.logger.WithField("listen", listen).Info("Stopping Metrics server")
}

func (pmd *PrometheusMetricsDriver) InitRoutes(r *mux.Router) {}

func (pmd *PrometheusMetricsDriver) ServeRequest(r *ServeRequest) {
	Requests.With(prometheus.Labels{
		"path":   r.ResolvedPath,
		"hash":   r.Hash,
		"client": r.RemoteAddr,
		"host":   r.Host,
	}).Observe(float64(r.Duration))
}
