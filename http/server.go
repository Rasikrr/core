package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/Rasikrr/core/api"
	"github.com/Rasikrr/core/config"
	"github.com/Rasikrr/core/log"
	"github.com/Rasikrr/core/metrics"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

const (
	readTimeout  = time.Minute
	writeTimeout = time.Minute
	idleTimeout  = 3 * time.Minute
)

type Server struct {
	name        string
	port        string
	host        string
	srv         *http.Server
	middlewares []Middleware
	router      *chi.Mux
}

func NewServer(
	_ context.Context,
	cfg config.HTTPConfig,
	metricsCfg config.Metrics,
	httpMetrics metrics.HTTPMetrics,
) *Server {
	router := chi.NewRouter()

	srv := &Server{
		name: cfg.Name,
		port: cfg.Port,
		host: cfg.Host,
		srv: &http.Server{
			Addr:         address(cfg.Host, cfg.Port),
			Handler:      router,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			IdleTimeout:  idleTimeout,
		},
		router: router,
	}
	if metricsCfg.Enabled {
		srv.WithMiddlewares(metrics.NewHTTPMetricsMiddleware(httpMetrics))
	}
	srv.WithMiddlewares(NewCORSMiddleware())
	srv.WithMiddlewares(NewRecoverMiddleware())
	srv.registerMiddlewares()
	addHealthRoute(router)
	return srv
}

func (s *Server) WithControllers(controllers ...Controller) {
	for _, c := range controllers {
		c.Init(s.router)
	}
}

func (s *Server) WithMiddlewares(middlewares ...Middleware) {
	s.middlewares = append(s.middlewares, middlewares...)
}

func (s *Server) Start(ctx context.Context) error {
	log.Infof(ctx, "starting %s http server on %s", s.name, address(s.host, s.port))
	if err := s.srv.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}

func (s *Server) registerMiddlewares() {
	for _, m := range s.middlewares {
		s.router.Use(m.Handle)
	}
	// use default chi middlewares
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)
}

func (s *Server) Close(ctx context.Context) error {
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Infof(ctx, "HTTP server shutdown error: %v", err)
		return fmt.Errorf("HTTP server shutdown error: %w", err)
	}
	log.Infof(ctx, "%s HTTP server closed", s.name)
	return nil
}

func address(host, port string) string {
	return host + ":" + port
}

func addHealthRoute(router *chi.Mux) {
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Info(r.Context(), "health check")
		api.SendData(w, map[string]string{
			"status": "OK",
		}, http.StatusOK)
	})
}
