package http

import (
	"github.com/Rasikrr/core/sentry"
	sentryhttp "github.com/getsentry/sentry-go/http"
)

func (s *Server) setupSentryMiddleware() {
	if !sentry.Enabled() {
		return
	}
	sentryMiddleware := sentryhttp.New(sentryhttp.Options{
		Repanic:         true, // ✅ Передаем panic в RecoverMiddleware
		WaitForDelivery: false,
	})

	s.router.Use(sentryMiddleware.Handle)
}
