package metrics

import (
	"context"
	"time"

	"beryju.org/imagik/pkg/config"
	"github.com/gorilla/mux"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type InfluxDBMetricsDriver struct {
	client influxdb2.Client
	api    api.WriteAPIBlocking
}

func (imd *InfluxDBMetricsDriver) Init() {
	imd.client = influxdb2.NewClient(config.C.MetricsInfluxDBConfig.URL, config.C.MetricsInfluxDBConfig.Token)
	imd.api = imd.client.WriteAPIBlocking(config.C.MetricsInfluxDBConfig.Org, config.C.MetricsInfluxDBConfig.Bucket)
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
	imd.api.WritePoint(context.Background(), p)
}
