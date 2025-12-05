package http

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	name        = "metrics"
	defaultHost = "0.0.0.0"
)

func NewMetricsServer(_ context.Context, port string) *Server {
	router := chi.NewRouter()

	srv := &Server{
		name: name,
		port: port,
		host: defaultHost,
		srv: &http.Server{
			Addr:         address(defaultHost, port),
			Handler:      router,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			IdleTimeout:  idleTimeout,
		},
		router: router,
	}
	srv.WithMiddlewares(NewRecoverMiddleware())
	srv.registerDefaultMiddlewares()

	router.Handle("/metrics", promhttp.Handler())

	return srv
}
