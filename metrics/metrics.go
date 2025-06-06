package metrics

import (
	"context"
)

type Metricer interface {
	HTTP() HTTPMetrics
	GRPCServer() GRPCServerMetrics
	GRPCClient() GRPCClientMetrics
}

type PrometheusMetricer struct {
	http       HTTPMetrics
	grpcServer GRPCServerMetrics
	grpcClient GRPCClientMetrics
}

func NewPrometheusMetricer(_ context.Context, namespace string) Metricer {
	return &PrometheusMetricer{
		http:       newHTTPMetrics(namespace),
		grpcServer: NewGRPCServerMetrics(namespace),
		grpcClient: NewGRPCClientMetrics(namespace),
	}
}

func (m *PrometheusMetricer) HTTP() HTTPMetrics {
	return m.http
}

func (m *PrometheusMetricer) GRPCServer() GRPCServerMetrics {
	return m.grpcServer
}

func (m *PrometheusMetricer) GRPCClient() GRPCClientMetrics {
	return m.grpcClient
}
