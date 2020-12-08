package main

import (
	"os"
	"time"

	"github.com/BeryJu/gopyazo/cmd"
	"github.com/getsentry/sentry-go"
)

func main() {
	if dsn, enabled := os.LookupEnv("SENTRY_DSN"); enabled {
		sentry.Init(sentry.ClientOptions{
			Dsn:              dsn,
			AttachStacktrace: true,
			TracesSampleRate: 1,
		})
		defer sentry.Flush(time.Second * 5)
		defer sentry.Recover()
	}
	cmd.Execute()
}
