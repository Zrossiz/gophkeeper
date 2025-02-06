// Package router provides the routing structure for user authentication and registration.
package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// UserRouter defines routes related to user authentication and registration.
type UserRouter struct {
	handler UserHandler // Handler for processing user-related requests.
}

// UserHandler defines the methods required for user authentication and registration.
type UserHandler interface {
	// Login handles user authentication requests.
	Login(rw http.ResponseWriter, r *http.Request)

	// Registration handles new user registration requests.
	Registration(rw http.ResponseWriter, r *http.Request)
}

// NewUserRouter creates a new instance of UserRouter.
//
// Parameters:
//   - h UserHandler: The handler implementing user authentication and registration logic.
//
// Returns:
//   - *UserRouter: A new instance of UserRouter.
func NewUserRouter(h UserHandler) *UserRouter {
	return &UserRouter{handler: h}
}

// RegisterRoutes registers user-related routes for authentication and registration.
//
// Parameters:
//   - r chi.Router: The router instance where user routes will be registered.
func (u *UserRouter) RegisterRoutes(r chi.Router) {
	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", u.handler.Registration) // Endpoint for user registration.
		r.Post("/login", u.handler.Login)           // Endpoint for user authentication.
	})
}
