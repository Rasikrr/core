package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	coreError "github.com/Rasikrr/core/errors"
)

type QueryParametersGetter interface {
	GetQueryParameters(r *http.Request) error
}

type ParametersGetter interface {
	GetParameters(r *http.Request) error
}

// nolint: errcheck
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

// nolint: errcheck
func SendError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	var (
		errorResp ErrorResponse
		coreErr   *coreError.CoreError
	)

	if errors.As(err, &coreErr) {
		errorResp.Error = coreErr.Message
		errorResp.StatusCode = coreErr.Code
	} else {
		coreErr = coreError.GRPCErrorToHTTP(err)
		errorResp.Error = coreErr.Message
		errorResp.StatusCode = coreErr.Code
	}

	bb, _ := json.Marshal(errorResp)
	w.WriteHeader(errorResp.StatusCode)
	w.Write(bb)
}

func GetData(r *http.Request, data interface{}) error {
	queryParams, ok := data.(QueryParametersGetter)
	if ok {
		if err := queryParams.GetQueryParameters(r); err != nil {
			return err
		}
		return nil
	}
	params, ok := data.(ParametersGetter)
	if ok {
		if err := params.GetParameters(r); err != nil {
			return err
		}
		return nil
	}
	defer r.Body.Close()
	bb, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if len(bb) > 0 {
		unmarshaller, ok := data.(json.Unmarshaler)
		if ok {
			return unmarshaller.UnmarshalJSON(bb)
		}
	}
	return nil
}
