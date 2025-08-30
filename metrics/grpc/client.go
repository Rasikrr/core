package grpc

import (
	"context"
	"time"

	"github.com/Rasikrr/core/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type ClientMetrics struct {
	reqs     *prometheus.CounterVec
	duration *prometheus.HistogramVec
	errors   *prometheus.CounterVec
}

func NewGRPCClientMetrics(m metrics.Metricer, clientName string) *ClientMetrics {
	ss := m.Subsystem("grpc_client")
	return &ClientMetrics{
		reqs:     ss.CounterVec(clientName+"requests_total", "Total gRPC client requests.", []string{"type", "method"}),
		duration: ss.HistogramVec(clientName+"request_duration_seconds", "gRPC client duration.", []string{"type", "method"}, prometheus.DefBuckets),
		errors:   ss.CounterVec(clientName+"errors_total", "Total gRPC client errors.", []string{"type", "method", "code"}),
	}
}

func (g *ClientMetrics) UnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption,
	) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		code := status.Code(err).String()

		g.reqs.WithLabelValues("unary", method).Inc()
		g.duration.WithLabelValues("unary", method).Observe(time.Since(start).Seconds())
		if code != "OK" {
			g.errors.WithLabelValues("unary", method, code).Inc()
		}
		return err
	}
}
