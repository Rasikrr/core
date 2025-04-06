package http

import (
	"github.com/go-chi/chi/v5"
)

type Controller interface {
	Init(router *chi.Mux)
}
