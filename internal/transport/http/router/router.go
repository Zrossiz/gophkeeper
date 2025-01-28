package router

type Router struct {
	Card     CardRouter
	User     UserRouter
	Binary   BinaryRouter
	LogoPass LogoPassRouter
}

func New() *Router {
	r := &Router{
		Card:     *NewCardRouter(),
		User:     *NewUserRouter(),
		Binary:   *NewBinaryRouter(),
		LogoPass: *NewLogoPassRouter(),
	}

	return r
}
