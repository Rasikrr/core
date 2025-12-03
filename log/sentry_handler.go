package log

import (
	"context"
	"log/slog"

	sentrySDK "github.com/getsentry/sentry-go"
	sentryslog "github.com/getsentry/sentry-go/slog"
)

func sentryHandler() slog.Handler {
	handler := sentryslog.Option{
		EventLevel: []slog.Level{slog.LevelError},
		LogLevel:   []slog.Level{}, // Это сбор логов, не нужен (используем loki или аналоги)
	}.NewSentryHandler(context.Background())

	breadCrumbHandler := newBreadcrumbHandler(handler)
	return breadCrumbHandler
}

// breadcrumbHandler adds breadcrumbs to Sentry Hub for all log levels
type breadcrumbHandler struct {
	next slog.Handler
}

func newBreadcrumbHandler(next slog.Handler) *breadcrumbHandler {
	return &breadcrumbHandler{next: next}
}

func (h *breadcrumbHandler) Enabled(ctx context.Context, level slog.Level) bool {
	// Always enabled to capture breadcrumbs for all levels
	return true
}

func (h *breadcrumbHandler) Handle(ctx context.Context, record slog.Record) error {
	// Only add breadcrumbs if we have a Hub in context (i.e., within an HTTP request)
	// This prevents breadcrumbs from startup logs polluting request-specific breadcrumbs
	hub := sentrySDK.GetHubFromContext(ctx)
	if hub != nil {
		// Collect attributes
		data := make(map[string]interface{})
		record.Attrs(func(attr slog.Attr) bool {
			data[attr.Key] = attr.Value.Any()
			return true
		})

		// Add breadcrumb to request-specific Hub
		hub.AddBreadcrumb(&sentrySDK.Breadcrumb{
			Type:      "default",
			Category:  "log",
			Message:   record.Message,
			Level:     convertSlogLevelToSentry(record.Level),
			Data:      data,
			Timestamp: record.Time,
		}, nil)
	}

	// Pass to next handler only if it's enabled for this level
	// (sentryslog only handles errors, but breadcrumbs are captured for all)
	if h.next.Enabled(ctx, record.Level) {
		return h.next.Handle(ctx, record)
	}
	return nil
}

func (h *breadcrumbHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &breadcrumbHandler{next: h.next.WithAttrs(attrs)}
}

func (h *breadcrumbHandler) WithGroup(name string) slog.Handler {
	return &breadcrumbHandler{next: h.next.WithGroup(name)}
}

func convertSlogLevelToSentry(level slog.Level) sentrySDK.Level {
	switch level {
	case slog.LevelDebug:
		return sentrySDK.LevelDebug
	case slog.LevelInfo:
		return sentrySDK.LevelInfo
	case slog.LevelWarn:
		return sentrySDK.LevelWarning
	case slog.LevelError:
		return sentrySDK.LevelError
	case LevelFatal:
		return sentrySDK.LevelFatal
	default:
		return sentrySDK.LevelInfo
	}
}
