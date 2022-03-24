package main

import (
	"fmt"
	"os"
	"time"

	"beryju.org/imagik/cmd"
	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
)

var buildCommit string

func main() {
	dsn := "https://bc5df9f742f04f62bb1e4b37b227aa62@sentry.beryju.org/7"
	if edsn, enabled := os.LookupEnv("SENTRY_DSN"); enabled {
		dsn = edsn
	}
	env := "default-env"
	if eenv, enabled := os.LookupEnv("SENTRY_ENVIRONMENT"); enabled {
		env = eenv
	}
	l := log.WithField("component", "imagik.root.sentry")
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		AttachStacktrace: true,
		TracesSampleRate: 1,
		Environment:      env,
		Debug:            true,
		DebugWriter:      l.Writer(),
		Release:          fmt.Sprintf("imagik@%s", buildCommit),
	})
	if err != nil {
		log.WithError(err).Warning("failed to init sentry")
	}
	log.WithField("commit", buildCommit).Info("imagik starting.")
	defer sentry.Flush(time.Second * 5)
	defer sentry.Recover()
	cmd.Execute()
}
