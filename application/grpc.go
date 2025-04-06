package application

import (
	"context"
	coreGrpc "github.com/Rasikrr/core/grpc"
	"log"
)

// nolint: unparam
func (a *App) initGRPC(_ context.Context) error {
	if !a.config.GRPC.Required {
		return nil
	}
	a.grpcServer = coreGrpc.NewServer(a.Config().GRPCConfig())
	log.Println("grpc initialized")
	a.starters.Add(a.grpcServer)
	a.closers.Add(a.grpcServer)
	return nil
}
