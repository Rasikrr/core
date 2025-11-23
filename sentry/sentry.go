package sentry

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Rasikrr/core/enum"
	"github.com/Rasikrr/core/version"
	sentrySDK "github.com/getsentry/sentry-go"
)

var (
	enabled atomic.Bool
	once    sync.Once
)

func Init(config Config, env enum.Environment) error {
	var initErr error
	once.Do(func() {
		initErr = sentrySDK.Init(sentrySDK.ClientOptions{
			Dsn:              config.DSN,
			Environment:      env.String(),
			Release:          version.GetVersion(),
			SampleRate:       config.SampleRate,
			TracesSampleRate: config.TracesSampleRate,
			Debug:            config.Debug,
			BeforeSend: func(event *sentrySDK.Event, _ *sentrySDK.EventHint) *sentrySDK.Event {
				return event
			},
		})
		if initErr == nil {
			enabled.Store(true)
		}
	})
	if initErr != nil {
		return fmt.Errorf("sentry.Init: %v", initErr)
	}
	return nil
}

// CaptureEvent отправляет событие в Sentry
func CaptureEvent(ctx context.Context, level slog.Level, msg string, attrs map[string]any) {
	if !Enabled() {
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
	if !Enabled() {
		return
	}
	sentrySDK.CaptureException(err)
}

// CaptureMessage отправляет сообщение в Sentry
func CaptureMessage(msg string) {
	if !Enabled() {
		return
	}
	sentrySDK.CaptureMessage(msg)
}

// Flush дожидается отправки всех событий (полезно перед завершением приложения)
func Flush(timeout time.Duration) bool {
	if !Enabled() {
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

func Enabled() bool {
	return enabled.Load()
}
