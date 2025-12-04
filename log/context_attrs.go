package log

import (
	"context"

	coreCtx "github.com/Rasikrr/core/context"

	"log/slog"
)

func getAttrsFromCtx(ctx context.Context) []slog.Attr {
	var attrs []slog.Attr
	if reqID, ok := coreCtx.RequestID(ctx); ok {
		attrs = append(attrs, slog.String(string(coreCtx.CtxKeyRequestID), reqID))
	}
	if userID, ok := coreCtx.UserID(ctx); ok {
		attrs = append(attrs, slog.String(string(coreCtx.CtxKeyUserID), userID))
	}
	return attrs
}
