package api

import (
	"encoding/json"
	"errors"
	coreError "github.com/Rasikrr/core/errors"
	"io"
	"net/http"
)

type ErrorResponse struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

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
			SendError(w, http.StatusInternalServerError, err)
			return
		}
		w.Write(bb)
		return
	}
	bb, err := json.Marshal(data)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}
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

// nolint: errcheck
func SendError(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	r := ErrorResponse{
		Error:  err.Error(),
		Status: statusCode,
	}
	var coreErr *coreError.CoreError

	if errors.As(err, &coreErr) {
		r.Error = coreErr.Message
		r.Status = coreErr.Code
	}
	bb, _ := json.Marshal(r)
	w.WriteHeader(r.Status)
	w.Write(bb)
}
