package log

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"time"

	"github.com/Rasikrr/core/sentry"
)

const (
	FatalString             = "FATAL"
	SentryString            = "SENTRY"
	LevelFatal   slog.Level = 12
	LevelSentry  slog.Level = 11
)

type Logger interface {
	Info(ctx context.Context, msg string, attrs ...Attr)
	Infof(ctx context.Context, format string, a ...any)

	Warn(ctx context.Context, msg string, attrs ...Attr)
	Warnf(ctx context.Context, format string, a ...any)

	Error(ctx context.Context, msg string, attrs ...Attr)
	Errorf(ctx context.Context, format string, a ...any)

	Debug(ctx context.Context, msg string, attrs ...Attr)
	Debugf(ctx context.Context, format string, a ...any)

	Fatal(ctx context.Context, msg string, attrs ...Attr)
	Fatalf(ctx context.Context, format string, a ...any)

	Sentry(ctx context.Context, msg string, attrs ...Attr)
	Sentryf(ctx context.Context, format string, a ...any)

	With(attrs ...Attr) Logger
}

type slogWrapper struct {
	base *slog.Logger
}

func (l *slogWrapper) log(ctx context.Context, level slog.Level, msg string, attrs []Attr) {
	if !l.base.Enabled(ctx, level) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(4, pcs[:])

	record := slog.NewRecord(
		time.Now(),
		level,
		msg,
		pcs[0],
	)

	for _, attr := range getAttrsFromCtx(ctx) {
		record.AddAttrs(attr)
	}
	for _, attr := range convertAttrsToSlog(attrs) {
		record.AddAttrs(attr)
	}

	_ = l.base.Handler().Handle(ctx, record)
}

func (l *slogWrapper) Info(ctx context.Context, msg string, attrs ...Attr) {
	l.log(ctx, slog.LevelInfo, msg, attrs)
}

func (l *slogWrapper) Infof(ctx context.Context, format string, a ...any) {
	l.log(ctx, slog.LevelInfo, fmt.Sprintf(format, a...), nil)
}

func (l *slogWrapper) Warn(ctx context.Context, msg string, attrs ...Attr) {
	l.log(ctx, slog.LevelWarn, msg, attrs)
}

func (l *slogWrapper) Warnf(ctx context.Context, format string, a ...any) {
	l.log(ctx, slog.LevelWarn, fmt.Sprintf(format, a...), nil)
}

func (l *slogWrapper) Error(ctx context.Context, msg string, attrs ...Attr) {
	l.log(ctx, slog.LevelError, msg, attrs)
}

func (l *slogWrapper) Errorf(ctx context.Context, format string, a ...any) {
	l.log(ctx, slog.LevelError, fmt.Sprintf(format, a...), nil)
}

func (l *slogWrapper) Debug(ctx context.Context, msg string, attrs ...Attr) {
	l.log(ctx, slog.LevelDebug, msg, attrs)
}

func (l *slogWrapper) Debugf(ctx context.Context, format string, a ...any) {
	l.log(ctx, slog.LevelDebug, fmt.Sprintf(format, a...), nil)
}

func (l *slogWrapper) Fatal(ctx context.Context, msg string, attrs ...Attr) {
	l.log(ctx, LevelFatal, msg, attrs)
	os.Exit(1)
}

func (l *slogWrapper) Fatalf(ctx context.Context, format string, a ...any) {
	l.log(ctx, LevelFatal, fmt.Sprintf(format, a...), nil)
	os.Exit(1)
}

func (l *slogWrapper) Sentry(ctx context.Context, msg string, attrs ...Attr) {
	l.log(ctx, LevelSentry, msg, attrs)
	sentry.CaptureEvent(ctx, LevelSentry, msg, attrsToMap(attrs))
}

func (l *slogWrapper) Sentryf(ctx context.Context, format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	l.log(ctx, LevelSentry, msg, nil)
	sentry.CaptureEvent(ctx, LevelSentry, msg, nil)
}

// attrsToMap конвертирует []Attr в map[string]any для Sentry
func attrsToMap(attrs []Attr) map[string]any {
	m := make(map[string]any, len(attrs))
	for _, attr := range attrs {
		m[attr.Key] = attr.Value.Any()
	}
	return m
}

func (l *slogWrapper) With(attrs ...Attr) Logger {
	return &slogWrapper{base: l.base.With(convertAttrsToAny(attrs)...)}
}
