package http

import (
	"net/http"

	coreCtx "github.com/Rasikrr/core/context"
	"github.com/Rasikrr/core/tracing"
	"github.com/google/uuid"
	"github.com/riandyrn/otelchi"
	"go.opentelemetry.io/otel/trace"
)

func (s *Server) setupTracingMiddleware() {
	if tracing.Enabled() {
		s.router.Use(otelchi.Middleware(s.name))
	}
	s.router.Use(traceCtxMiddleware)
}

func traceCtxMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var traceID string

		// Приоритет 1: traceID из OpenTelemetry span (если трейсинг включен)
		sc := trace.SpanContextFromContext(ctx)
		if sc.HasTraceID() {
			traceID = sc.TraceID().String()
		}

		// Приоритет 2: traceID из заголовка запроса
		if traceID == "" {
			traceID = r.Header.Get(TraceIDHeader)
		}

		// Приоритет 3: генерируем новый UUID
		if traceID == "" {
			traceID = uuid.New().String()
		}

		ctx = coreCtx.WithTraceID(ctx, traceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
