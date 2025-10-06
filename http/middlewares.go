package http

import (
	"fmt"
	"net/http"

	"github.com/Rasikrr/core/log"
	"github.com/go-chi/cors"
)

type Middleware interface {
	Handle(next http.Handler) http.Handler
}

type CORSMiddleware struct {
	handler func(next http.Handler) http.Handler
}

func NewCORSMiddleware(options cors.Options) *CORSMiddleware {
	return &CORSMiddleware{
		handler: cors.Handler(options),
	}
}

func (c *CORSMiddleware) Handle(next http.Handler) http.Handler {
	return c.handler(next)
}

type RecoverMiddleware struct{}

func NewRecoverMiddleware() Middleware {
	return &RecoverMiddleware{}
}

func (m *RecoverMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		defer func() {
			if err := recover(); err != nil {
				log.Debugf(ctx, "panic while handling request: %v", err)
				http.Error(w, fmt.Sprintf("panic: %v", err), http.StatusInternalServerError)
				return
			}
		}()
		next.ServeHTTP(w, r)
	})
}
