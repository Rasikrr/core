package log

import (
	"context"
	"github.com/Rasikrr/core/enum"
	"log/slog"
	"os"
	"sync"
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

func Init(env enum.Environment) {
	once.Do(func() {
		var handler slog.Handler
		switch env {
		case enum.EnvironmentProd:
			handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level:     slog.LevelInfo,
				AddSource: false,
			})
		default:
			handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level:     slog.LevelDebug,
				AddSource: true,
			})
		}
		defaultLogger = &slogWrapper{base: slog.New(handler)}
	})
}
