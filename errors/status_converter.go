package errors

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

var (
	clientClosedRequest = 499

	gRPCCodesToHTTP = map[codes.Code]int{
		codes.OK:                 http.StatusOK,
		codes.Canceled:           clientClosedRequest,
		codes.Unknown:            http.StatusInternalServerError,
		codes.InvalidArgument:    http.StatusBadRequest,
		codes.DeadlineExceeded:   http.StatusRequestTimeout,
		codes.NotFound:           http.StatusNotFound,
		codes.AlreadyExists:      http.StatusConflict,
		codes.PermissionDenied:   http.StatusForbidden,
		codes.ResourceExhausted:  http.StatusTooManyRequests,
		codes.FailedPrecondition: http.StatusBadRequest,
		codes.Aborted:            http.StatusConflict,
		codes.OutOfRange:         http.StatusBadRequest,
		codes.Unimplemented:      http.StatusNotImplemented,
		codes.Internal:           http.StatusInternalServerError,
		codes.Unavailable:        http.StatusServiceUnavailable,
		codes.DataLoss:           http.StatusInternalServerError,
		codes.Unauthenticated:    http.StatusUnauthorized,
	}
)

func GRPCErrorToHTTP(err error) *CoreError {
	s, ok := status.FromError(err)
	if !ok {
		err = fmt.Errorf("unknown error: %w", err)
		return NewError(err.Error(), http.StatusInternalServerError)
	}
	code, ok := gRPCCodesToHTTP[s.Code()]
	if !ok {
		err = fmt.Errorf("unknown gRPC code: %s", s.Code())
		return NewError(err.Error(), http.StatusInternalServerError)
	}
	return NewError(s.Message(), code)
}
