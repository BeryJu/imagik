package main

import (
	"os"
	"time"

	"beryju.org/imagik/cmd"
	"github.com/getsentry/sentry-go"
)

func main() {
	dsn := "https://759fc52c64334ea0a2978460a6469fd0@sentry.beryju.org/15"
	if edsn, enabled := os.LookupEnv("SENTRY_DSN"); enabled {
		dsn = edsn
	}
	env := "default-env"
	if eenv, enabled := os.LookupEnv("SENTRY_ENVIRONMENT"); enabled {
		env = eenv
	}
	sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		AttachStacktrace: true,
		TracesSampleRate: 1,
		Environment:      env,
	})
	defer sentry.Flush(time.Second * 5)
	defer sentry.Recover()
	cmd.Execute()
}
