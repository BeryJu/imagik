package metrics

import (
	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

type NullMetricsDriver struct {
	logger *log.Entry
}

func (nmd *NullMetricsDriver) Init() {
	nmd.logger = log.WithField("component", "metrics")
}

func (nmd *NullMetricsDriver) InitRoutes(r *mux.Router) {}

func (nmd *NullMetricsDriver) ServeRequest(r *ServeRequest) {
	nmd.logger.WithFields(log.Fields{
		"Path":      r.ResolvedPath,
		"Hash":      r.Hash,
		"Client":    r.RemoteAddr,
		"UserAgent": r.UserAgent(),
	}).Info(r.URL.Path)
}
