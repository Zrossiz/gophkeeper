package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type LogoPassRouter struct {
	h LogoPassHandler
	m Middleware
}

type LogoPassHandler interface {
	GetAll(rw http.ResponseWriter, r *http.Request)
	Update(rw http.ResponseWriter, r *http.Request)
	Create(rw http.ResponseWriter, r *http.Request)
}

func NewLogoPassRouter(
	h LogoPassHandler,
	m Middleware,
) *LogoPassRouter {
	return &LogoPassRouter{
		h: h,
		m: m,
	}
}

func (c *LogoPassRouter) RegisterRoutes(r chi.Router) {
	r.Route("/api/logo-pass", func(r chi.Router) {
		r.With(c.m.Auth).Post("/", c.h.Create)
		r.With(c.m.Auth).Get("/{userID}", c.h.GetAll)
		r.With(c.m.Auth).Put("/{logoPassID}", c.h.Update)
	})
}
