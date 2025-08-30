package grpc

import (
	"context"
	"time"

	"github.com/Rasikrr/core/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type ServerMetrics struct {
	reqs     *prometheus.CounterVec
	duration *prometheus.HistogramVec
	errors   *prometheus.CounterVec
}

func NewGRPCServerMetrics(m metrics.Metricer) *ServerMetrics {
	ss := m.Subsystem("grpc_server")
	return &ServerMetrics{
		reqs:     ss.CounterVec("requests_total", "Total gRPC server requests.", []string{"type", "method"}),
		duration: ss.HistogramVec("request_duration_seconds", "gRPC server duration.", []string{"type", "method"}, prometheus.DefBuckets),
		errors:   ss.CounterVec("errors_total", "Total gRPC server errors.", []string{"type", "method", "code"}),
	}
}

func (g *ServerMetrics) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()

		resp, err := handler(ctx, req)
		code := status.Code(err).String()
		g.reqs.WithLabelValues("unary", info.FullMethod).Inc()
		g.duration.WithLabelValues("unary", info.FullMethod).Observe(time.Since(start).Seconds())

		if code != "OK" {
			g.errors.WithLabelValues("unary", info.FullMethod, code).Inc()
		}

		return resp, err
	}
}
