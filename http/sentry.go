package http

import (
	"time"

	"github.com/Rasikrr/core/sentry"
	sentryhttp "github.com/getsentry/sentry-go/http"
)

func (s *Server) setupSentryMiddleware() {
	if !sentry.Enabled() {
		return
	}
	sentryMiddleware := sentryhttp.New(sentryhttp.Options{
		Repanic:         true,  // Передаем panic дальше после захвата
		WaitForDelivery: false, // Не блокируем ответ (для production)
		Timeout:         2 * time.Second,
	})

	s.router.Use(sentryMiddleware.Handle)
}
