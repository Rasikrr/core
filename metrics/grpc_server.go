package metrics

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"time"
)

type GRPCServerMetrics interface {
	IncRequest(typ, method string)
	ObserveDuration(typ, method string, seconds float64)
	IncError(typ, method, code string)
}

type grpcServerMetrics struct {
	reqs     *prometheus.CounterVec
	duration *prometheus.HistogramVec
	errors   *prometheus.CounterVec
}

func NewGRPCServerMetrics(namespace string) GRPCServerMetrics {
	m := &grpcServerMetrics{
		reqs: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "grpc",
				Name:      "requests_total",
				Help:      "Total number of gRPC requests received.",
			},
			[]string{"type", "method"},
		),
		duration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: "grpc",
				Name:      "request_duration_seconds",
				Help:      "Duration of gRPC requests.",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"type", "method"},
		),
		errors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "grpc",
				Name:      "errors_total",
				Help:      "Total number of gRPC request errors.",
			},
			[]string{"type", "method", "code"},
		),
	}

	prometheus.MustRegister(m.reqs, m.duration, m.errors)
	return m
}

func (m *grpcServerMetrics) IncRequest(typ, method string) {
	m.reqs.WithLabelValues(typ, method).Inc()
}

func (m *grpcServerMetrics) ObserveDuration(typ, method string, seconds float64) {
	m.duration.WithLabelValues(typ, method).Observe(seconds)
}

func (m *grpcServerMetrics) IncError(typ, method, code string) {
	m.errors.WithLabelValues(typ, method, code).Inc()
}

type MyInt int

func Test() MyInt {
	return 34
}

func UnaryServerInterceptor(m GRPCServerMetrics) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()

		resp, err := handler(ctx, req)
		code := status.Code(err).String()

		m.IncRequest("unary", info.FullMethod)
		m.ObserveDuration("unary", info.FullMethod, time.Since(start).Seconds())

		if code != "OK" {
			m.IncError("unary", info.FullMethod, code)
		}

		return resp, err
	}
}
