package sentrylog

import (
	"fmt"

	"gitlab.boomerangapp.ir/back/utils/pkg/conv"
)

const RequestID = "request_id"

//SentryWriter for set log output
type SentryWriter struct{}

func (sn *SentryWriter) Write(p []byte) (n int, err error) {
	if !isActive {
		fmt.Printf("%s", p)
		return
	}

	Info(conv.B2S(p))
	Recover()
	return len(p), nil
}
