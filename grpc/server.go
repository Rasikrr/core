package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/Rasikrr/core/log"
	"github.com/Rasikrr/core/tracing"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

const (
	TCP = "tcp"
	UDP = "udp"
)

type Server struct {
	host   string
	port   int
	server *grpc.Server
}

func NewServer(
	cfg Config,
) *Server {
	return &Server{
		host:   cfg.Host,
		port:   cfg.Port,
		server: newGrpcServer(),
	}
}

func (s *Server) Start(ctx context.Context) error {
	lis, err := net.Listen(TCP, s.addr(s.host))
	if err != nil {
		return err
	}
	log.Info(ctx, "starting grpc server")
	if err := s.server.Serve(lis); err != nil {
		if errors.Is(err, grpc.ErrServerStopped) {
			return nil
		}
		return err
	}
	return nil
}

func (s *Server) Close(ctx context.Context) error {
	s.server.GracefulStop()
	log.Info(ctx, "grpc server closed")
	return nil
}

func (s *Server) Srv() *grpc.Server {
	return s.server
}

func (s *Server) addr(host string) string {
	return fmt.Sprintf("%s:%d", host, s.port)
}

func newGrpcServer() *grpc.Server {
	initGRPCMetrics()

	unaryInterceptors := []grpc.UnaryServerInterceptor{
		unaryPanicRecoveryInterceptor,
		UnaryServerSentryInterceptor,
		metrics.UnaryServer(),
	}

	if tracing.Enabled() {
		unaryInterceptors = append(unaryInterceptors, UnaryServerTraceInterceptor)
	}

	streamInterceptors := []grpc.StreamServerInterceptor{
		streamPanicRecoveryInterceptor,
		StreamServerSentryInterceptor,
		metrics.StreamServer(),
	}

	if tracing.Enabled() {
		streamInterceptors = append(streamInterceptors, StreamServerTraceInterceptor)
	}

	unary := grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			unaryInterceptors...,
		),
	)
	stream := grpc.StreamInterceptor(
		grpc_middleware.ChainStreamServer(
			streamInterceptors...,
		),
	)

	opts := []grpc.ServerOption{unary, stream}
	if tracing.Enabled() {
		opts = append(opts, tracingServerInterceptor())
	}

	return grpc.NewServer(opts...)
}
