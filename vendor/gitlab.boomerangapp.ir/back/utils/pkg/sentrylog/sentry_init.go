package sentrylog

import (
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	"gitlab.boomerangapp.ir/back/utils/pkg/env"
)

var isActive bool

//Init init sentry connection
func Init(Dsn string, Version string) {
	if Dsn == "" {
		Dsn = env.Str("SENTRY_DSN", "")
	}
	if Dsn == "" {
		isActive = false
		return
	}
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              Dsn,
		Debug:            env.Bool("DEBUG", false),
		ServerName:       env.Str("SERVER_NAME", "Not Set"),
		Release:          Version,
		TracesSampleRate: env.Float("SENTRY_SAMPLERATE", 1.0),
	})
	if err != nil {
		log.Print(err)
		return
	}

	//set active flag
	isActive = true
	sentry.Flush(time.Second * 7)
}

//Recover set `defer Recover()` to the main function
func Recover() {
	sentry.Flush(time.Second * 7)
	sentry.Recover()
}
