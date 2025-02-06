// Package router defines the HTTP routing structure for handling different entities such as users, cards, notes, binary data, and logo passwords.
package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Router holds all the sub-routers responsible for handling different API routes.
type Router struct {
	Card     CardRouter     // Routes for card-related operations.
	User     UserRouter     // Routes for user-related operations.
	Binary   BinaryRouter   // Routes for binary data operations.
	LogoPass LogoPassRouter // Routes for logo password operations.
	Note     NoteRouter     // Routes for note-related operations.
}

// Handler contains the handlers required for processing API requests.
type Handler struct {
	User     UserHandler     // Handler for user-related operations.
	Card     CardHandler     // Handler for card-related operations.
	Binary   BinaryHandler   // Handler for binary data operations.
	LogoPass LogoPassHandler // Handler for logo password operations.
	Note     NoteHandler     // Handler for note-related operations.
}

// Middleware defines an interface for handling authentication middleware.
type Middleware interface {
	// Auth applies authentication middleware to a given HTTP handler.
	Auth(next http.Handler) http.Handler
}

// New initializes a new HTTP router with registered routes for handling API requests.
//
// Parameters:
//   - h Handler: A struct containing the handlers for various entities.
//   - m Middleware: An implementation of authentication middleware.
//
// Returns:
//   - *chi.Mux: A new router with registered routes.
func New(h Handler, m Middleware) *chi.Mux {
	r := chi.NewRouter()

	// Initialize and assign routers for different functionalities.
	router := &Router{
		Card:     *NewCardRouter(h.Card, m),
		User:     *NewUserRouter(h.User),
		Binary:   *NewBinaryRouter(h.Binary, m),
		LogoPass: *NewLogoPassRouter(h.LogoPass, m),
		Note:     *NewNoteRouter(h.Note, m),
	}

	// Register routes for each module.
	router.User.RegisterRoutes(r)
	router.Card.RegisterRoutes(r)
	router.LogoPass.RegisterRoutes(r)
	router.Binary.RegisterRoutes(r)
	router.Note.RegisterRoutes(r)

	// Register Swagger documentation handler.
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	return r
}
