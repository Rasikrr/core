package log

import (
	"context"
	"log/slog"
)

type ctxKey string

const (
	CtxKeyRequestID ctxKey = "request_id"
	CtxKeyUserID    ctxKey = "user_id"
)

func getAttrsFromCtx(ctx context.Context) []slog.Attr {
	var attrs []slog.Attr
	if reqID, ok := ctx.Value(CtxKeyRequestID).(string); ok {
		attrs = append(attrs, slog.String("request_id", reqID))
	}
	if userID, ok := ctx.Value(CtxKeyUserID).(string); ok {
		attrs = append(attrs, slog.String("user_id", userID))
	}
	return attrs
}
