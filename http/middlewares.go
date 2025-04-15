package http

import (
	"fmt"
	"github.com/Rasikrr/core/log"
	"net/http"
)

type Middleware interface {
	Handle(next http.Handler) http.Handler
}

type CORSMiddleware struct{}

func NewCORSMiddleware() Middleware {
	return &CORSMiddleware{}
}

func (m *CORSMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
		} else {
			next.ServeHTTP(w, r)
		}
	})
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
