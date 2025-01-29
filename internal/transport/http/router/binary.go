package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type BinaryRouter struct {
	h BinaryHandler
}

type BinaryHandler interface {
	GetAll(rw http.ResponseWriter, r *http.Request)
	Update(rw http.ResponseWriter, r *http.Request)
	Create(rw http.ResponseWriter, r *http.Request)
}

func NewBinaryRouter(h BinaryHandler) *BinaryRouter {
	return &BinaryRouter{
		h: h,
	}
}

func (b *BinaryRouter) RegisterRoutes(r chi.Router) {
	r.Route("/binary", func(r chi.Router) {
		r.Post("/", b.h.Create)
		r.Get("/", b.h.GetAll)
		r.Put("/{id}", b.h.Update)
	})
}
