package http

import (
	"context"
	"net/http"

	"github.com/Rasikrr/core/config"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	name        = "metrics"
	defaultHost = "0.0.0.0"
)

func NewMetricsServer(_ context.Context, cfg config.Metrics) *Server {
	router := chi.NewRouter()

	srv := &Server{
		name: name,
		port: cfg.Prometheus.Port,
		host: defaultHost,
		srv: &http.Server{
			Addr:         address(defaultHost, cfg.Prometheus.Port),
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
