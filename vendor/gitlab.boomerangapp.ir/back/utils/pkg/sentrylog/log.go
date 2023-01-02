package sentrylog

import (
	"encoding/json"
	"log"
	"os"

	"github.com/getsentry/sentry-go"
	"gitlab.boomerangapp.ir/back/utils/pkg/conv"
)

//Info log info messages
func Info(data interface{}, tags ...string) {
	if !isActive {
		log.Print(data)
		return
	}

	sentry.WithScope(func(scope *sentry.Scope) {
		if len(tags) > 0 {
			scope.SetTag(RequestID, tags[0])
		}
		scope.SetLevel(sentry.LevelInfo)
		switch data.(type) {
		case string:
			sentry.CaptureMessage(data.(string))
		default:
			message, _ := json.Marshal(data)
			sentry.CaptureMessage(conv.B2S(message))
		}
	})
}

//Warning set warning log
func Warning(err error, tags ...string) {
	Log(err, sentry.LevelWarning, tags...)
}

//Fatal set fatal logs
func Fatal(err error, tags ...string) {
	if !isActive {
		log.Fatal(err)
		return
	}

	Log(err, sentry.LevelFatal, tags...)

	//recover and send data
	Recover()

	os.Exit(1)
}

//Log log errors
func Log(err error, level sentry.Level, tags ...string) {

	if !isActive {
		log.Print(err)
		return
	}

	sentry.WithScope(func(scope *sentry.Scope) {
		if len(tags) > 0 {
			scope.SetTag(RequestID, tags[0])
		}

		scope.SetLevel(level)
		sentry.CaptureException(err)
	})
}
