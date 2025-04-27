package log

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"time"
)

const (
	FatalString            = "FATAL"
	LevelFatal  slog.Level = 12
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

func (l *slogWrapper) With(attrs ...Attr) Logger {
	return &slogWrapper{base: l.base.With(convertAttrsToAny(attrs)...)}
}
