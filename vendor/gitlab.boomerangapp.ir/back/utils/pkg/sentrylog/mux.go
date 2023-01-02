package sentrylog

import (
	"context"
	"errors"
	"net/http"

	"github.com/getsentry/sentry-go"
	"gitlab.boomerangapp.ir/back/utils/consts"
)

//SentryMiddleware is the sentry helper middleware to use in mux routers
func SentryMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rID := r.Header.Get(consts.HTTPRequestID)
		if rID == "" {
			Warning(errors.New("request without id"))
		}

		hub := sentry.CurrentHub().Clone()
		hub.Scope().SetTag(RequestID, rID)

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), valuesKey, hub)))
	})
}

//GetHubFromCtx retrieves attached *sentry.Hub instance from context
func GetHubFromCtx(ctx context.Context) *sentry.Hub {

	if hub, ok := ctx.Value(valuesKey).(*sentry.Hub); ok {
		return hub
	}

	return nil
}
