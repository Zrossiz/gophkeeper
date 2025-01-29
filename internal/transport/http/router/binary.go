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
	Update(rw http.ResponseWriter, r *http.Request)
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
		r.Post("/", b.h.Create)
		r.Get("/", b.h.GetAll)
		r.Put("/{binaryID}", b.h.Update)
	})
}
