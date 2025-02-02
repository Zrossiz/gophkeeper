package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Router struct {
	Card     CardRouter
	User     UserRouter
	Binary   BinaryRouter
	LogoPass LogoPassRouter
	Note     NoteRouter
}

type Handler struct {
	User     UserHandler
	Card     CardHandler
	Binary   BinaryHandler
	LogoPass LogoPassHandler
	Note     NoteHandler
}

type Middleware interface {
	Auth(next http.Handler) http.Handler
}

func New(
	h Handler,
	m Middleware,
) *chi.Mux {
	r := chi.NewRouter()

	router := &Router{
		Card:     *NewCardRouter(h.Card, m),
		User:     *NewUserRouter(h.User),
		Binary:   *NewBinaryRouter(h.Binary, m),
		LogoPass: *NewLogoPassRouter(h.LogoPass, m),
		Note:     *NewNoteRouter(h.Note, m),
	}

	router.User.RegisterRoutes(r)
	router.Card.RegisterRoutes(r)
	router.LogoPass.RegisterRoutes(r)
	router.Binary.RegisterRoutes(r)
	router.Note.RegisterRoutes(r)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	return r
}
