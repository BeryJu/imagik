package metrics

import (
	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

type NullMetricsDriver struct {
	logger *log.Entry
}

func (nmd *NullMetricsDriver) Init() {
	nmd.logger = log.WithField("component", "imagik.drivers.metrics.null")
}

func (nmd *NullMetricsDriver) InitRoutes(r *mux.Router) {}

func (nmd *NullMetricsDriver) ServeRequest(r *ServeRequest) {
}
