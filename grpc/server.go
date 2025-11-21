package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"runtime/debug"

	"github.com/Rasikrr/core/log"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		metrics.UnaryServer(),
	}
	streamInterceptors := []grpc.StreamServerInterceptor{
		streamPanicRecoveryInterceptor,
		metrics.StreamServer(),
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
	return grpc.NewServer(unary, stream)
}

func streamPanicRecoveryInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	_ *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	ctx := ss.Context()
	defer func() {
		if r := recover(); r != nil {
			log.Debugf(ctx, "Recovered from panic in stream: %v\n%s", r, debug.Stack())
		}
	}()
	return handler(srv, ss)
}

func unaryPanicRecoveryInterceptor(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Debugf(ctx, "Recovered from panic in unary: %v\n%s", r, debug.Stack())
			err = status.Errorf(codes.Internal, "internal server error")
		}
	}()
	return handler(ctx, req)
}
