package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Rasikrr/core/log"
)

type Middleware interface {
	Handle(next http.Handler) http.Handler
}

type CORSMiddleware struct {
	origins []string
	methods []string
	headers []string
	creds   bool
}

func NewCORSMiddleware() *CORSMiddleware {
	return &CORSMiddleware{}
}

func (c *CORSMiddleware) WithOrigins(origins ...string) *CORSMiddleware {
	c.origins = append(c.origins, origins...)
	return c
}

func (c *CORSMiddleware) WithMethods(methods ...string) *CORSMiddleware {
	c.methods = append(c.methods, methods...)
	return c
}

func (c *CORSMiddleware) WithHeaders(headers ...string) *CORSMiddleware {
	c.headers = append(c.headers, headers...)
	return c
}

func (c *CORSMiddleware) WithCredentials(creds bool) *CORSMiddleware {
	c.creds = creds
	return c
}

func (c *CORSMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", strings.Join(c.origins, ", "))
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(c.methods, ", "))
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(c.headers, ", "))
		w.Header().Set("Access-Control-Allow-Credentials", fmt.Sprintf("%v", c.creds))

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
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
