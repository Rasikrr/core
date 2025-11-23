package http

import (
	"net/http"

	"github.com/Rasikrr/core/api"
	coreError "github.com/Rasikrr/core/errors"
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
		wrappedWriter := &responseWriter{
			ResponseWriter: w,
			statusCode:     0,
			written:        false,
		}

		defer func() {
			if err := recover(); err != nil {
				log.Errorf(ctx, "panic while handling request: %v", err)
				if !wrappedWriter.written {
					api.SendError(w, coreError.ErrInternalServerError)
				} else {
					log.Warnf(ctx, "cannot set status 500: response already written with status %d", wrappedWriter.statusCode)
				}
				return
			}
		}()

		next.ServeHTTP(wrappedWriter, r)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	if !rw.written {
		rw.statusCode = statusCode
		rw.written = true
		rw.ResponseWriter.WriteHeader(statusCode)
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.statusCode = http.StatusOK
		rw.written = true
	}
	return rw.ResponseWriter.Write(b)
}
