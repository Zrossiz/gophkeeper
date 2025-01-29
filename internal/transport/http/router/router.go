package router

import "github.com/go-chi/chi/v5"

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

func New(h Handler) *chi.Mux {
	r := chi.NewRouter()

	router := &Router{
		Card:     *NewCardRouter(h.Card),
		User:     *NewUserRouter(h.User),
		Binary:   *NewBinaryRouter(h.Binary),
		LogoPass: *NewLogoPassRouter(h.LogoPass),
	}

	router.User.RegisterRoutes(r)
	router.Card.RegisterRoutes(r)
	router.LogoPass.RegisterRoutes(r)
	router.Binary.RegisterRoutes(r)

	return r
}
