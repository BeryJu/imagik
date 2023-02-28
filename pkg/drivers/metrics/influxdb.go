package metrics

import (
	"context"
	"time"

	"beryju.io/imagik/pkg/config"
	"github.com/gorilla/mux"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	log "github.com/sirupsen/logrus"
)

type InfluxDBMetricsDriver struct {
	client influxdb2.Client
	api    api.WriteAPIBlocking
	logger *log.Entry
}

func (imd *InfluxDBMetricsDriver) Init() {
	imd.client = influxdb2.NewClient(config.C.MetricsInfluxDBConfig.URL, config.C.MetricsInfluxDBConfig.Token)
	imd.api = imd.client.WriteAPIBlocking(config.C.MetricsInfluxDBConfig.Org, config.C.MetricsInfluxDBConfig.Bucket)
	imd.logger = log.WithField("component", "imagik.drivers.metrics.influx")
}

func (imd *InfluxDBMetricsDriver) InitRoutes(r *mux.Router) {}

func (imd *InfluxDBMetricsDriver) ServeRequest(r *ServeRequest) {
	p := influxdb2.NewPointWithMeasurement("imagik_serve").
		AddTag("path", r.ResolvedPath).
		AddTag("hash", r.Hash).
		AddTag("client", r.RemoteAddr).
		AddTag("user-agent", r.UserAgent()).
		AddField("request", 1).
		SetTime(time.Now())
	err := imd.api.WritePoint(context.Background(), p)
	if err != nil {
		imd.logger.WithError(err).Warning("failed to write points")
	}
}
