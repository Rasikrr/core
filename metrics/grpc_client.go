package metrics

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"time"
)

type GRPCClientMetrics interface {
	IncRequest(typ, method string)
	ObserveDuration(typ, method string, seconds float64)
	IncError(typ, method, code string)
}

type grpcClientMetrics struct {
	reqs     *prometheus.CounterVec
	duration *prometheus.HistogramVec
	errors   *prometheus.CounterVec
}

func NewGRPCClientMetrics(namespace string) GRPCClientMetrics {
	m := &grpcClientMetrics{
		reqs: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "grpc_client",
				Name:      "requests_total",
				Help:      "Total number of gRPC client requests.",
			},
			[]string{"type", "method"},
		),
		duration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: "grpc_client",
				Name:      "request_duration_seconds",
				Help:      "Duration of gRPC client requests.",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"type", "method"},
		),
		errors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "grpc_client",
				Name:      "errors_total",
				Help:      "Total number of gRPC client errors.",
			},
			[]string{"type", "method", "code"},
		),
	}

	prometheus.MustRegister(m.reqs, m.duration, m.errors)
	return m
}

func (m *grpcClientMetrics) IncRequest(typ, method string) {
	m.reqs.WithLabelValues(typ, method).Inc()
}

func (m *grpcClientMetrics) ObserveDuration(typ, method string, seconds float64) {
	m.duration.WithLabelValues(typ, method).Observe(seconds)
}

func (m *grpcClientMetrics) IncError(typ, method, code string) {
	m.errors.WithLabelValues(typ, method, code).Inc()
}

func UnaryClientInterceptor(m GRPCClientMetrics) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		code := status.Code(err).String()

		m.IncRequest("unary", method)
		m.ObserveDuration("unary", method, time.Since(start).Seconds())

		if code != "OK" {
			m.IncError("unary", method, code)
		}

		return err
	}
}
