package application

import (
	"context"
	coreGrpc "github.com/Rasikrr/core/grpc"
	"github.com/Rasikrr/core/log"
	"github.com/Rasikrr/core/metrics"
)

// nolint: unparam
func (a *App) initGRPC(ctx context.Context) error {
	if !a.config.GRPC.Required {
		return nil
	}

	var grpcMetrics metrics.GRPCServerMetrics
	if a.metrics != nil {
		grpcMetrics = a.metrics.GRPCServer()
	}
	a.grpcServer = coreGrpc.NewServer(
		a.Config().GRPCConfig(),
		a.Config().MetricsConfig(),
		grpcMetrics,
	)

	log.Info(ctx, "grpc initialized")

	a.starters.Add(a.grpcServer)
	a.closers.Add(a.grpcServer)

	return nil
}
