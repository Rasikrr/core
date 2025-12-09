// nolint: errcheck
package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	coreCtx "github.com/Rasikrr/core/context"
)

func SendData(ctx context.Context, w http.ResponseWriter, data interface{}, statusCode int) {
	traceID, ok := coreCtx.TraceID(ctx)
	if ok {
		w.Header().Set(TraceIDHeader, traceID)
	}
	w.Header().Set(ContentTypeHeader, "application/json")
	w.WriteHeader(statusCode)

	marshaller, ok := data.(json.Marshaler)
	if ok {
		bb, err := marshaller.MarshalJSON()
		if err != nil {
			SendError(ctx, w, err)
			return
		}
		w.Write(bb)
		return
	}
	bb, err := json.Marshal(data)
	if err != nil {
		SendError(ctx, w, err)
		return
	}
	w.Write(bb)
}

func SendError(ctx context.Context, w http.ResponseWriter, err error) {
	traceID, ok := coreCtx.TraceID(ctx)
	if ok {
		w.Header().Set(TraceIDHeader, traceID)
	}
	w.Header().Set("Content-Type", "application/json")

	var (
		errorResp ErrorResponse
		httpError *Error
	)

	if errors.As(err, &httpError) {
		errorResp.Message = httpError.Message
		errorResp.Code = httpError.Code
	} else {
		errorResp.Message = err.Error()
		errorResp.Code = http.StatusInternalServerError
	}

	w.WriteHeader(errorResp.Code)

	bb, _ := json.Marshal(errorResp)
	w.Write(bb)
}
