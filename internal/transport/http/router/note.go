package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type NoteRouter struct {
	h LogoPassHandler
	m Middleware
}

type NoteHandler interface {
	GetAll(rw http.ResponseWriter, r *http.Request)
	Update(rw http.ResponseWriter, r *http.Request)
	Create(rw http.ResponseWriter, r *http.Request)
}

func NewNoteRouter(
	h NoteHandler,
	m Middleware,
) *NoteRouter {
	return &NoteRouter{
		h: h,
		m: m,
	}
}

func (n *NoteRouter) RegisterRoutes(r chi.Router) {
	r.Route("/api/note", func(r chi.Router) {
		r.With(n.m.Auth).Post("/", n.h.Create)
		r.With(n.m.Auth).Get("/{userID}", n.h.GetAll)
		r.With(n.m.Auth).Put("/{noteID}", n.h.Update)
	})
}
