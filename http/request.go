package http

import (
	"encoding/json"
	"io"
	"net/http"
)

type QueryParametersGetter interface {
	GetQueryParameters(r *http.Request) error
}

type ParametersGetter interface {
	GetParameters(r *http.Request) error
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
		unmarshaler, ok := data.(json.Unmarshaler)
		if ok {
			return unmarshaler.UnmarshalJSON(bb)
		}
		return json.Unmarshal(bb, data)
	}
	return nil
}
