package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type LogoPassRouter struct {
	h LogoPassHandler
}

type LogoPassHandler interface {
	GetAll(rw http.ResponseWriter, r *http.Request)
	Update(rw http.ResponseWriter, r *http.Request)
	Create(rw http.ResponseWriter, r *http.Request)
}

func NewLogoPassRouter(h LogoPassHandler) *LogoPassRouter {
	return &LogoPassRouter{
		h: h,
	}
}

func (c *LogoPassRouter) RegisterRoutes(r chi.Router) {
	r.Route("/api/logo-pass", func(r chi.Router) {
		r.Post("/", c.h.Create)
		r.Get("/", c.h.GetAll)
		r.Put("/{logoPassID}", c.h.Update)
	})
}
