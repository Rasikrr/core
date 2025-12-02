package log

import (
	"context"
	"log/slog"
	"os"
	"sync"

	"github.com/Rasikrr/core/enum"
	"github.com/Rasikrr/core/sentry"
	"github.com/cockroachdb/errors"
	sentryslog "github.com/getsentry/sentry-go/slog"
	slogmulti "github.com/samber/slog-multi"
)

var (
	defaultLogger Logger = &slogWrapper{base: slog.New(slog.NewTextHandler(os.Stdout, nil))}
	once          sync.Once
)

func Info(ctx context.Context, msg string, args ...Attr) {
	Default().Info(ctx, msg, args...)
}

func Infof(ctx context.Context, format string, a ...any) {
	Default().Infof(ctx, format, a...)
}

func Warn(ctx context.Context, msg string, args ...Attr) {
	Default().Warn(ctx, msg, args...)
}

func Warnf(ctx context.Context, format string, a ...any) {
	Default().Warnf(ctx, format, a...)
}

func Error(ctx context.Context, msg string, args ...Attr) {
	Default().Error(ctx, msg, args...)
}

func Errorf(ctx context.Context, format string, a ...any) {
	Default().Errorf(ctx, format, a...)
}

func Debug(ctx context.Context, msg string, args ...Attr) {
	Default().Debug(ctx, msg, args...)
}

func Debugf(ctx context.Context, format string, a ...any) {
	Default().Debugf(ctx, format, a...)
}

func Fatal(ctx context.Context, msg string, args ...Attr) {
	Default().Fatal(ctx, msg, args...)
}

func Fatalf(ctx context.Context, format string, args ...any) {
	Default().Fatalf(ctx, format, args...)
}

func With(args ...Attr) Logger {
	return defaultLogger.With(args...)
}

func Default() Logger {
	return defaultLogger
}

func Init(env enum.Environment, level enum.LogLevel, addSource bool) error {
	lvl := level.ToSlogLevel()
	opts := &slog.HandlerOptions{
		Level:       lvl,
		AddSource:   addSource,
		ReplaceAttr: replaceLevelAttr,
	}
	var onceErr error
	once.Do(func() {
		handlers := make([]slog.Handler, 0, 2)
		if sentry.Enabled() {
			sentryHandler, err := getSentryHandler()
			if err != nil {
				onceErr = errors.Wrap(err, "logger init failed")
				return
			}
			handlers = append(handlers, sentryHandler)
		}

		switch env {
		case enum.EnvironmentProd:
			handlers = append(handlers, slog.NewJSONHandler(os.Stdout, opts))
		default:
			handlers = append(handlers, slog.NewTextHandler(os.Stdout, opts))
		}
		handler := slogmulti.Fanout(handlers...)

		defaultLogger = &slogWrapper{base: slog.New(handler)}
	})
	return onceErr
}

// nolint: gocritic
func replaceLevelAttr(_ []string, a slog.Attr) slog.Attr {
	if a.Key == slog.LevelKey {
		level := a.Value.Any().(slog.Level)
		if level == LevelFatal {
			return slog.String(slog.LevelKey, FatalString)
		}
	}
	return a
}

func getSentryHandler() (slog.Handler, error) {
	return sentryslog.Option{
		EventLevel: []slog.Level{slog.LevelError},
		LogLevel:   []slog.Level{slog.LevelWarn, slog.LevelInfo},
	}.NewSentryHandler(context.Background()), nil
}
