package transports

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/luciodesimone/go-kit-example/endpoints"
)

func NewHTTPHandler(e endpoints.Set, logger log.Logger) http.Handler {
	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	m := http.NewServeMux()
	m.Handle("/healthz", httptransport.NewServer(
		e.HealthEndpoint,
		decodeHTTPHealthCheckRequest,
		encodeHTTPGenericResponse,
		opts...,
	))

	return m
}

func decodeHTTPHealthCheckRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return r, nil
}

func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		errorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(errToCode(err))
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}

type errorWrapper struct {
	Error string `json:"error"`
}

func errToCode(err error) int {
	switch err {
	// case TODO complete when errors defined
	}
	return http.StatusInternalServerError
}
