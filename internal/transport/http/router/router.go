package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	Card     CardRouter
	User     UserRouter
	Binary   BinaryRouter
	LogoPass LogoPassRouter
}

type Handler struct {
	User     UserHandler
	Card     CardHandler
	Binary   BinaryHandler
	LogoPass LogoPassHandler
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
	}

	router.User.RegisterRoutes(r)
	router.Card.RegisterRoutes(r)
	router.LogoPass.RegisterRoutes(r)
	router.Binary.RegisterRoutes(r)

	return r
}
