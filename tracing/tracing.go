package tracing

import (
	"context"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace/noop"
)

var (
	once    sync.Once
	enabled atomic.Bool
)

func Init(ctx context.Context, cfg Config, appName string) error {
	if !cfg.Enabled {
		otel.SetTracerProvider(noop.NewTracerProvider())
		return nil
	}

	// 1. Создаём OTLP gRPC экспортёр (по умолчанию на localhost:4317, можно переопределить)
	var err error

	once.Do(func() {
		// 1. Настраиваем propagator для передачи trace context между сервисами
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		))

		var exporter *otlptrace.Exporter
		exporter, err = otlptracegrpc.New(ctx,
			otlptracegrpc.WithEndpoint(cfg.DSN),
			otlptracegrpc.WithInsecure(),
		)
		if err != nil {
			return
		}

		// 2. Создаём ресурс (service name + метаданные)
		var res *resource.Resource
		res, err = resource.Merge(
			resource.Default(),
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(appName),
			),
		)
		if err != nil {
			return
		}

		// 3. Создаём TracerProvider с batch-процессором
		tp := tracesdk.NewTracerProvider(
			tracesdk.WithBatcher(exporter),
			tracesdk.WithResource(res),
		)

		// 4. Устанавливаем глобальный провайдер (optional)
		otel.SetTracerProvider(tp)
		enabled.Store(true)
	})
	return err
}

func Enabled() bool {
	return enabled.Load()
}

func GetTracer(name string) Tracer {
	return otel.Tracer(name)
}
