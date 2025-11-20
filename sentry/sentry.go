package sentry

import (
	"context"
	"log/slog"
	"time"

	"github.com/Rasikrr/core/enum"
	sentrySDK "github.com/getsentry/sentry-go"
)

var enabled bool

// Init инициализирует Sentry SDK
func Init(dsn string, environment enum.Environment, opts ...Option) error {
	config := &Config{
		SampleRate:       1.0,
		TracesSampleRate: 0.1,
		Debug:            false,
	}

	for _, opt := range opts {
		opt(config)
	}

	err := sentrySDK.Init(sentrySDK.ClientOptions{
		Dsn:              dsn,
		Environment:      environment.String(),
		SampleRate:       config.SampleRate,
		TracesSampleRate: config.TracesSampleRate,
		Debug:            config.Debug,
		BeforeSend: func(event *sentrySDK.Event, hint *sentrySDK.EventHint) *sentrySDK.Event {
			// Можно добавить фильтрацию или модификацию событий
			return event
		},
	})

	if err == nil {
		enabled = true
	}

	return err
}

// CaptureEvent отправляет событие в Sentry
func CaptureEvent(ctx context.Context, level slog.Level, msg string, attrs map[string]any) {
	if !enabled {
		return
	}

	event := sentrySDK.NewEvent()
	event.Message = msg
	event.Level = convertLevel(level)
	event.Timestamp = time.Now()

	// Добавляем контекст из ctx
	if requestID, ok := ctx.Value("request_id").(string); ok {
		event.Tags["request_id"] = requestID
	}
	if userID, ok := ctx.Value("user_id").(string); ok {
		event.User = sentrySDK.User{ID: userID}
	}

	// Конвертируем атрибуты в tags и extra
	for key, value := range attrs {
		// Ошибки добавляем как exceptions
		if key == "error" {
			if err, ok := value.(error); ok {
				event.Exception = []sentrySDK.Exception{{
					Value: err.Error(),
					Type:  "error",
				}}
				continue
			}
		}

		// Простые типы в tags, сложные в extra
		switch v := value.(type) {
		case string, int, int64, bool, float64:
			event.Tags[key] = toString(v)
		default:
			event.Extra[key] = v
		}
	}

	sentrySDK.CaptureEvent(event)
}

// CaptureException отправляет ошибку в Sentry
func CaptureException(err error) {
	if !enabled {
		return
	}
	sentrySDK.CaptureException(err)
}

// CaptureMessage отправляет сообщение в Sentry
func CaptureMessage(msg string) {
	if !enabled {
		return
	}
	sentrySDK.CaptureMessage(msg)
}

// Flush дожидается отправки всех событий (полезно перед завершением приложения)
func Flush(timeout time.Duration) bool {
	if !enabled {
		return true
	}
	return sentrySDK.Flush(timeout)
}

// convertLevel конвертирует slog.Level в sentrySDK.Level
func convertLevel(level slog.Level) sentrySDK.Level {
	switch {
	case level >= 12: // Fatal
		return sentrySDK.LevelFatal
	case level >= 11: // Sentry
		return sentrySDK.LevelError
	case level >= slog.LevelError:
		return sentrySDK.LevelError
	case level >= slog.LevelWarn:
		return sentrySDK.LevelWarning
	case level >= slog.LevelInfo:
		return sentrySDK.LevelInfo
	default:
		return sentrySDK.LevelDebug
	}
}

// toString конвертирует значение в строку для tags
func toString(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case int:
		return slog.IntValue(val).String()
	case int64:
		return slog.Int64Value(val).String()
	case bool:
		return slog.BoolValue(val).String()
	case float64:
		return slog.Float64Value(val).String()
	default:
		return ""
	}
}
