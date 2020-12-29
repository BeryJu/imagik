package metrics

import (
	"fmt"
	"net/http"
	"os"

	"github.com/BeryJu/imagik/pkg/config"
	"github.com/BeryJu/imagik/pkg/drivers"
	"github.com/gorilla/mux"
)

type ServeRequest struct {
	http.Request
	Hash         string
	ResolvedPath string
}

func NewServeRequest(r *http.Request) *ServeRequest {
	return &ServeRequest{
		Request:      *r,
		Hash:         "",
		ResolvedPath: "",
	}
}

type MetricsDriver interface {
	drivers.HTTPDriver
	ServeRequest(r *ServeRequest)
}

func FromConfig(r *mux.Router) MetricsDriver {
	metricsDriverType := config.C.MetricsDriver
	var metricsDriver MetricsDriver
	switch metricsDriverType {
	case "null":
		metricsDriver = &NullMetricsDriver{}
	case "influxdb":
		metricsDriver = &InfluxDBMetricsDriver{}
	}
	if metricsDriver == nil {
		fmt.Printf("Could not configure metricsDriver '%s'", metricsDriverType)
		os.Exit(1)
	}
	metricsDriver.Init()
	metricsDriver.InitRoutes(r)
	return metricsDriver
}
