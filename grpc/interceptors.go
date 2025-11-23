package grpc

import (
	"context"
	"runtime/debug"

	"github.com/Rasikrr/core/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func streamPanicRecoveryInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	_ *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	ctx := ss.Context()
	defer func() {
		if r := recover(); r != nil {
			log.Warnf(ctx, "Recovered from panic in stream: %v\n%s", r, debug.Stack())
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
			log.Warnf(ctx, "Recovered from panic in unary: %v\n%s", r, debug.Stack())
			err = status.Errorf(codes.Internal, "internal server error")
		}
	}()
	return handler(ctx, req)
}
