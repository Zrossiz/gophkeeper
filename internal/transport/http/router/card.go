package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type CardRouter struct {
	h CardHandler
}

type CardHandler interface {
	GetAll(rw http.ResponseWriter, r *http.Request)
	Update(rw http.ResponseWriter, r *http.Request)
	Create(rw http.ResponseWriter, r *http.Request)
}

func NewCardRouter(h CardHandler) *CardRouter {
	return &CardRouter{
		h: h,
	}
}

func (c *CardRouter) RegisterRoutes(r chi.Router) {
	r.Route("/card", func(r chi.Router) {
		r.Post("/", c.h.Create)
		r.Get("/", c.h.GetAll)
		r.Put("/{id}", c.h.Update)
	})
}
