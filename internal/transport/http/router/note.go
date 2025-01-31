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
		r.Post("/", n.h.Create)
		r.Get("/", n.h.GetAll)
		r.Put("/{noteID}", n.h.Update)
	})
}
