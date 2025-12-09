package grpc

import (
	"context"

	coreCtx "github.com/Rasikrr/core/context"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

func tracingClientInterceptor() grpc.DialOption {
	return grpc.WithStatsHandler(otelgrpc.NewClientHandler())
}

func tracingServerInterceptor() grpc.ServerOption {
	return grpc.StatsHandler(otelgrpc.NewServerHandler())
}

func UnaryServerTraceInterceptor(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	sc := trace.SpanContextFromContext(ctx)
	if sc.HasTraceID() {
		ctx = coreCtx.WithTraceID(ctx, sc.TraceID().String())
	}

	return handler(ctx, req)
}

type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedStream) Context() context.Context {
	return w.ctx
}

func StreamServerTraceInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	_ *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	sc := trace.SpanContextFromContext(ss.Context())
	if sc.IsValid() && sc.HasTraceID() {
		newCtx := coreCtx.WithTraceID(ss.Context(), sc.TraceID().String())

		ss = &wrappedStream{
			ServerStream: ss,
			ctx:          newCtx,
		}
	}

	return handler(srv, ss)
}
