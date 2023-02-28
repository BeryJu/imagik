package metrics

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"beryju.io/imagik/pkg/config"
	"beryju.io/imagik/pkg/drivers"
	"github.com/gorilla/mux"
)

type ServeRequest struct {
	http.Request
	Hash         string
	ResolvedPath string
	Duration     time.Duration
}

func NewServeRequest(r *http.Request) *ServeRequest {
	return &ServeRequest{
		Request:      *r,
		Hash:         "",
		ResolvedPath: "",
	}
}

type Driver interface {
	drivers.HTTPDriver
	ServeRequest(r *ServeRequest)
}

func FromConfig(r *mux.Router) Driver {
	metricsDriverType := config.C.MetricsDriver
	var metricsDriver Driver
	switch metricsDriverType {
	case "null":
		metricsDriver = &NullMetricsDriver{}
	case "influxdb":
		metricsDriver = &InfluxDBMetricsDriver{}
	case "prometheus":
		metricsDriver = &PrometheusMetricsDriver{}
	}
	if metricsDriver == nil {
		fmt.Printf("Could not configure metricsDriver '%s'", metricsDriverType)
		os.Exit(1)
	}
	metricsDriver.Init()
	metricsDriver.InitRoutes(r)
	return metricsDriver
}
