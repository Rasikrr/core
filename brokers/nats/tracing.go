package nats

// OpenTelemetry трейсинг для NATS
//
// Этот файл реализует распространение trace context через NATS сообщения.
// Trace context передается через заголовки NATS сообщений в формате W3C Trace Context.
//
// Publisher автоматически инжектирует trace context из context.Context в заголовки сообщения.
// Subscriber автоматически извлекает trace context из заголовков и создает новый span
// для обработки каждого сообщения.
//
// Использование:
//
// 1. Publisher:
//   ctx := context.Background()
//   publisher.Publish(ctx, "subject", message) // trace context автоматически добавится в заголовки
//
// 2. Subscriber:
//   type MyHandler struct{}
//   func (h *MyHandler) Handle(ctx context.Context, msg *nats.Msg) error {
//       // ctx содержит trace context из заголовков сообщения
//       // можно создавать child spans используя этот context
//       return nil
//   }
//   func (h *MyHandler) Subject() string { return "my.subject" }

import (
	"context"
	"fmt"

	coreCtx "github.com/Rasikrr/core/context"
	"github.com/Rasikrr/core/tracing"
	"github.com/nats-io/nats.go"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	tracerName = "github.com/Rasikrr/core/brokers/nats"
)

// natsHeaderCarrier адаптер для NATS заголовков, реализующий TextMapCarrier
type natsHeaderCarrier struct {
	header Header
}

func (c natsHeaderCarrier) Get(key string) string {
	return c.header.Get(key)
}

func (c natsHeaderCarrier) Set(key, value string) {
	c.header.Set(key, value)
}

func (c natsHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(c.header))
	for k := range c.header {
		keys = append(keys, k)
	}
	return keys
}

// injectTraceContext инжектирует trace context в заголовки NATS сообщения
func injectTraceContext(ctx context.Context, msg *nats.Msg) {
	if !tracing.Enabled() {
		return
	}

	carrier := natsHeaderCarrier{header: msg.Header}
	otel.GetTextMapPropagator().Inject(ctx, carrier)
}

// extractTraceContext извлекает trace context из заголовков NATS сообщения
// и создает новый span для обработки сообщения
func extractTraceContext(ctx context.Context, msg *nats.Msg, spanName string) (context.Context, trace.Span) {
	if !tracing.Enabled() {
		return ctx, trace.SpanFromContext(ctx)
	}

	tracer := tracing.GetTracer(tracerName)
	carrier := natsHeaderCarrier{header: msg.Header}

	// Извлекаем trace context из заголовков
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	// Создаем новый span для обработки сообщения
	ctx, span := tracer.Start(ctx,
		spanName,
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(
			attribute.String("messaging.system", "nats"),
			attribute.String("messaging.topic", msg.Subject),
			attribute.Int("messaging.message_payload_size_bytes", len(msg.Data)),
		),
	)

	// Добавляем trace ID в кастомный контекст
	sc := trace.SpanContextFromContext(ctx)
	if sc.HasTraceID() {
		ctx = coreCtx.WithTraceID(ctx, sc.TraceID().String())
	}

	return ctx, span
}

// startPublishSpan создает span для публикации сообщения
func startPublishSpan(ctx context.Context, subject string) (context.Context, trace.Span) {
	if !tracing.Enabled() {
		return ctx, trace.SpanFromContext(ctx)
	}

	tracer := tracing.GetTracer(tracerName)
	ctx, span := tracer.Start(ctx,
		fmt.Sprintf("nats.publish %s", subject),
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(
			attribute.String("messaging.system", "nats"),
			attribute.String("messaging.topic", subject),
		),
	)

	return ctx, span
}

// recordSpanError записывает ошибку в span
func recordSpanError(span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
}
