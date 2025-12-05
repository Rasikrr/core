package context

import "context"

type ctxKey string

const (
	CtxKeyRequestID ctxKey = "request_id"
	CtxKeyUserID    ctxKey = "user_id"
)

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, CtxKeyRequestID, id)
}

func RequestID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(CtxKeyRequestID).(string)
	return v, ok
}

func WithUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, CtxKeyUserID, id)
}

func UserID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(CtxKeyUserID).(string)
	return v, ok
}
