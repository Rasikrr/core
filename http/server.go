package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/Rasikrr/core/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"time"
)

const (
	readTimeout  = time.Minute
	writeTimeout = time.Minute
	idleTimeout  = 3 * time.Minute
)

type Server struct {
	port        string
	host        string
	srv         *http.Server
	middlewares []Middleware
	router      *chi.Mux
}

func NewServer(_ context.Context, cfg config.HTTPConfig) *Server {
	router := chi.NewRouter()

	srv := &Server{
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
	srv.WithMiddlewares(NewCORSMiddleware())
	srv.WithMiddlewares(NewRecoverMiddleware())
	srv.registerMiddlewares()
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

func (s *Server) Start(_ context.Context) error {
	log.Println("starting http server")
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
		log.Printf("HTTP server shutdown error: %v", err)
		return fmt.Errorf("HTTP server shutdown error: %w", err)
	}
	log.Println("HTTP server closed")
	return nil
}

func address(host, port string) string {
	return host + ":" + port
}
