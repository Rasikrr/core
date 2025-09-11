package application

import (
	"context"

	coreGrpc "github.com/Rasikrr/core/grpc"
	"github.com/Rasikrr/core/log"
)

// nolint: unparam
func (a *App) initGRPC(ctx context.Context) error {
	if !a.config.GRPC.Required {
		return nil
	}

	a.grpcServer = coreGrpc.NewServer(
		a.Config().GRPCConfig(),
	)

	log.Info(ctx, "grpc initialized")

	a.starters.Add(a.grpcServer)
	a.closers.Add(a.grpcServer)

	return nil
}
