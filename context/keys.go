package context

import "context"

type ctxKey string

const (
	CtxKeyTraceID ctxKey = "trace_id"
	CtxKeyUserID  ctxKey = "user_id"
)

func WithTraceID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, CtxKeyTraceID, id)
}

func TraceID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(CtxKeyTraceID).(string)
	return v, ok
}

func WithUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, CtxKeyUserID, id)
}

func UserID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(CtxKeyUserID).(string)
	return v, ok
}
