package http

import (
	"net/http"
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
	// Clear breadcrumbs for each request to prevent accumulation
	s.router.Use(clearBreadcrumbsMiddleware, sentryMiddleware.Handle)
}

// clearBreadcrumbsMiddleware clears breadcrumbs at the start of each request
// to ensure only request-specific breadcrumbs are captured
func clearBreadcrumbsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if hub := sentry.GetHubFromContext(r.Context()); hub != nil {
			hub.Scope().ClearBreadcrumbs()
		}
		next.ServeHTTP(w, r)
	})
}
