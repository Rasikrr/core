package http

import (
	"net/http"

	coreCtx "github.com/Rasikrr/core/context"
	"github.com/Rasikrr/core/tracing"
	"github.com/riandyrn/otelchi"
	"go.opentelemetry.io/otel/trace"
)

func (s *Server) setupTracingMiddleware() {
	if !tracing.Enabled() {
		return
	}
	s.router.Use(otelchi.Middleware(s.name), traceCtxMiddleware)
}

func traceCtxMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		sc := trace.SpanContextFromContext(ctx)
		if sc.HasTraceID() {
			ctx = coreCtx.WithTraceID(ctx, sc.TraceID().String())
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
