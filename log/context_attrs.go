package log

import "context"

type ctxKey string

const (
	CtxKeyRequestID ctxKey = "request_id"
	CtxKeyUserID    ctxKey = "user_id"
)

func getAttrsFromCtx(ctx context.Context) []any {
	var attrs []any
	if reqID, ok := ctx.Value(CtxKeyRequestID).(string); ok {
		attrs = append(attrs, "request_id", reqID)
	}
	if userID, ok := ctx.Value(CtxKeyUserID).(string); ok {
		attrs = append(attrs, "user_id", userID)
	}
	return attrs
}
