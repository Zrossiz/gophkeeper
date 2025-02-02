package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type BinaryRouter struct {
	h BinaryHandler
	m Middleware
}

type BinaryHandler interface {
	GetAll(rw http.ResponseWriter, r *http.Request)
	Create(rw http.ResponseWriter, r *http.Request)
}

func NewBinaryRouter(
	h BinaryHandler,
	m Middleware,
) *BinaryRouter {
	return &BinaryRouter{
		h: h,
		m: m,
	}
}

func (b *BinaryRouter) RegisterRoutes(r chi.Router) {
	r.Route("/api/binary", func(r chi.Router) {
		r.With(b.m.Auth).Post("/", b.h.Create)
		r.With(b.m.Auth).Get("/{userID}", b.h.GetAll)
	})
}
