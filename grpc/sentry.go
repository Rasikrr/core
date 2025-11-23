package grpc

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/Rasikrr/core/sentry"
	sentrySDK "github.com/getsentry/sentry-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryServerSentryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	if !sentry.Enabled() {
		return handler(ctx, req)
	}

	hub := sentrySDK.CurrentHub().Clone()
	ctx = sentrySDK.SetHubOnContext(ctx, hub)

	service, method := split(info.FullMethod)
	hub.Scope().SetContext("grpc", map[string]interface{}{
		"type":    "unary",
		"service": service,
		"method":  method,
	})

	defer func() {
		if r := recover(); r != nil {
			hub.WithScope(func(scope *sentrySDK.Scope) {
				scope.SetLevel(sentrySDK.LevelFatal)
				scope.SetTag("grpc.type", "unary")
				scope.SetTag("grpc.service", service)
				scope.SetTag("grpc.method", method)
				scope.SetContext("panic", map[string]interface{}{
					"value":      fmt.Sprintf("%v", r),
					"stacktrace": string(debug.Stack()),
				})

				var panicErr error
				switch x := r.(type) {
				case error:
					panicErr = x
				case string:
					panicErr = fmt.Errorf("panic: %s", x)
				default:
					panicErr = fmt.Errorf("panic: %v", r)
				}

				hub.CaptureException(panicErr)
			})
			panic(r) // re-panic, чтобы другие интерцепторы могли обработать
		}
	}()

	resp, err = handler(ctx, req)

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			code := st.Code()
			if shouldReportToSentry(code) {
				hub.WithScope(func(scope *sentrySDK.Scope) {
					scope.SetLevel(sentryLevelFromGRPCCode(code))
					scope.SetTag("grpc.type", "unary")
					scope.SetTag("grpc.service", service)
					scope.SetTag("grpc.method", method)
					scope.SetTag("grpc.code", code.String())
					scope.SetContext("grpc_error", map[string]interface{}{
						"message": st.Message(),
						"code":    code.String(),
						"details": st.Details(),
					})
					hub.CaptureException(err)
				})
			}
		}
	}

	return resp, err
}

// nolint: ineffassign, staticcheck
func StreamServerSentryInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	if !sentry.Enabled() {
		return handler(srv, ss)
	}

	ctx := ss.Context()
	hub := sentrySDK.CurrentHub().Clone()
	ctx = sentrySDK.SetHubOnContext(ctx, hub)

	service, method := split(info.FullMethod)
	hub.Scope().SetContext("grpc", map[string]interface{}{
		"type":    "stream",
		"service": service,
		"method":  method,
	})

	defer func() {
		if r := recover(); r != nil {
			hub.WithScope(func(scope *sentrySDK.Scope) {
				scope.SetLevel(sentrySDK.LevelFatal)
				scope.SetTag("grpc.type", "stream")
				scope.SetTag("grpc.service", service)
				scope.SetTag("grpc.method", method)
				scope.SetContext("panic", map[string]interface{}{
					"value":      fmt.Sprintf("%v", r),
					"stacktrace": string(debug.Stack()),
				})

				var panicErr error
				switch x := r.(type) {
				case error:
					panicErr = x
				case string:
					panicErr = fmt.Errorf("panic: %s", x)
				default:
					panicErr = fmt.Errorf("panic: %v", r)
				}

				hub.CaptureException(panicErr)
			})
			panic(r) // re-panic, чтобы другие интерцепторы могли обработать
		}
	}()

	err := handler(srv, &sentryServerStream{ServerStream: ss, hub: hub, service: service, method: method})

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			code := st.Code()
			if shouldReportToSentry(code) {
				hub.WithScope(func(scope *sentrySDK.Scope) {
					scope.SetLevel(sentryLevelFromGRPCCode(code))
					scope.SetTag("grpc.type", "stream")
					scope.SetTag("grpc.service", service)
					scope.SetTag("grpc.method", method)
					scope.SetTag("grpc.code", code.String())
					scope.SetContext("grpc_error", map[string]interface{}{
						"message": st.Message(),
						"code":    code.String(),
						"details": st.Details(),
					})
					hub.CaptureException(err)
				})
			}
		}
	}

	return err
}

type sentryServerStream struct {
	grpc.ServerStream
	hub     *sentrySDK.Hub
	service string
	method  string
}

func (s *sentryServerStream) RecvMsg(m interface{}) error {
	err := s.ServerStream.RecvMsg(m)
	if err != nil && shouldReportStreamError(err) {
		s.hub.WithScope(func(scope *sentrySDK.Scope) {
			scope.SetTag("grpc.type", "stream")
			scope.SetTag("grpc.service", s.service)
			scope.SetTag("grpc.method", s.method)
			scope.SetTag("stream.operation", "recv")
			s.hub.CaptureException(err)
		})
	}
	return err
}

func (s *sentryServerStream) SendMsg(m interface{}) error {
	err := s.ServerStream.SendMsg(m)
	if err != nil && shouldReportStreamError(err) {
		s.hub.WithScope(func(scope *sentrySDK.Scope) {
			scope.SetTag("grpc.type", "stream")
			scope.SetTag("grpc.service", s.service)
			scope.SetTag("grpc.method", s.method)
			scope.SetTag("stream.operation", "send")
			s.hub.CaptureException(err)
		})
	}
	return err
}

func shouldReportToSentry(code codes.Code) bool {
	switch code {
	case codes.OK,
		codes.Canceled,
		codes.InvalidArgument,
		codes.NotFound,
		codes.AlreadyExists,
		codes.PermissionDenied,
		codes.Unauthenticated,
		codes.FailedPrecondition,
		codes.OutOfRange:
		return false
	case codes.Unknown,
		codes.DeadlineExceeded,
		codes.ResourceExhausted,
		codes.Aborted,
		codes.Unimplemented,
		codes.Internal,
		codes.Unavailable,
		codes.DataLoss:
		return true
	default:
		return false
	}
}

func shouldReportStreamError(err error) bool {
	if err == nil {
		return false
	}
	st, ok := status.FromError(err)
	if !ok {
		return true
	}
	return shouldReportToSentry(st.Code())
}

func sentryLevelFromGRPCCode(code codes.Code) sentrySDK.Level {
	switch code {
	case codes.Internal, codes.DataLoss:
		return sentrySDK.LevelError
	case codes.Unknown, codes.Unavailable:
		return sentrySDK.LevelError
	case codes.DeadlineExceeded, codes.ResourceExhausted:
		return sentrySDK.LevelWarning
	case codes.Aborted, codes.Unimplemented:
		return sentrySDK.LevelWarning
	default:
		return sentrySDK.LevelInfo
	}
}
