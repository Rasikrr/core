package log

import (
	"context"
	"fmt"
	"log/slog"
	"os"
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

func (l *slogWrapper) Info(ctx context.Context, msg string, attrs ...Attr) {
	anyAttrs := convertAttrs(attrs)
	l.base.With(getAttrsFromCtx(ctx)...).InfoContext(ctx, msg, anyAttrs...)
}

func (l *slogWrapper) Infof(ctx context.Context, format string, a ...any) {
	l.base.With(getAttrsFromCtx(ctx)...).InfoContext(ctx, fmt.Sprintf(format, a...))
}

func (l *slogWrapper) Warn(ctx context.Context, msg string, attrs ...Attr) {
	anyAttrs := convertAttrs(attrs)
	l.base.With(getAttrsFromCtx(ctx)...).WarnContext(ctx, msg, anyAttrs...)
}

func (l *slogWrapper) Warnf(ctx context.Context, format string, a ...any) {
	l.base.With(getAttrsFromCtx(ctx)...).WarnContext(ctx, fmt.Sprintf(format, a...))
}

func (l *slogWrapper) Error(ctx context.Context, msg string, attrs ...Attr) {
	anyAttrs := convertAttrs(attrs)
	l.base.With(getAttrsFromCtx(ctx)...).ErrorContext(ctx, msg, anyAttrs...)
}

func (l *slogWrapper) Errorf(ctx context.Context, format string, a ...any) {
	l.base.With(getAttrsFromCtx(ctx)...).ErrorContext(ctx, fmt.Sprintf(format, a...))
}

func (l *slogWrapper) Debug(ctx context.Context, msg string, attrs ...Attr) {
	anyAttrs := convertAttrs(attrs)
	l.base.With(getAttrsFromCtx(ctx)...).DebugContext(ctx, msg, anyAttrs...)
}

func (l *slogWrapper) Debugf(ctx context.Context, format string, attrs ...any) {
	l.base.With(getAttrsFromCtx(ctx)...).DebugContext(ctx, fmt.Sprintf(format, attrs...))
}

func (l *slogWrapper) Fatal(ctx context.Context, msg string, attrs ...Attr) {
	l.Error(ctx, msg, attrs...)
	os.Exit(1)
}

func (l *slogWrapper) Fatalf(ctx context.Context, format string, a ...any) {
	l.Errorf(ctx, format, a...)
	os.Exit(1)
}

func (l *slogWrapper) With(attrs ...Attr) Logger {
	anyAttrs := convertAttrs(attrs)
	return &slogWrapper{base: l.base.With(anyAttrs...)}
}
