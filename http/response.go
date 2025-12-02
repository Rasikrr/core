// nolint: errcheck
package http

import (
	"encoding/json"
	"errors"
	"net/http"
)

func SendData(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	marshaller, ok := data.(json.Marshaler)
	if ok {
		bb, err := marshaller.MarshalJSON()
		if err != nil {
			SendError(w, err)
			return
		}
		w.Write(bb)
		return
	}
	bb, err := json.Marshal(data)
	if err != nil {
		SendError(w, err)
		return
	}
	w.Write(bb)
}

func SendError(w http.ResponseWriter, err error) {
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
