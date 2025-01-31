package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type CardRouter struct {
	h CardHandler
	m Middleware
}

type CardHandler interface {
	GetAll(rw http.ResponseWriter, r *http.Request)
	Update(rw http.ResponseWriter, r *http.Request)
	Create(rw http.ResponseWriter, r *http.Request)
}

func NewCardRouter(
	h CardHandler,
	m Middleware,
) *CardRouter {
	return &CardRouter{
		h: h,
		m: m,
	}
}

func (c *CardRouter) RegisterRoutes(r chi.Router) {
	r.Route("/api/card", func(r chi.Router) {
		r.With(c.m.Auth).Post("/", c.h.Create)
		r.With(c.m.Auth).Get("/", c.h.GetAll)
		r.With(c.m.Auth).Put("/{cardID}", c.h.Update)
	})
}
