// Package router defines the HTTP routing structure for handling logo-password-related requests.
package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// LogoPassRouter provides route registration for logo-password-related HTTP handlers.
type LogoPassRouter struct {
	h LogoPassHandler // Handler for logo-password operations.
	m Middleware      // Middleware for authentication and request processing.
}

// LogoPassHandler defines the interface for handling logo-password data requests.
type LogoPassHandler interface {
	// GetAll retrieves all stored logo-password information for a specific user.
	GetAll(rw http.ResponseWriter, r *http.Request)

	// Update modifies an existing logo-password entry.
	Update(rw http.ResponseWriter, r *http.Request)

	// Create adds a new logo-password entry to the storage.
	Create(rw http.ResponseWriter, r *http.Request)
}

// NewLogoPassRouter initializes a new LogoPassRouter instance.
//
// Parameters:
//   - h LogoPassHandler: The handler for logo-password operations.
//   - m Middleware: Middleware for handling authentication and authorization.
//
// Returns:
//   - *LogoPassRouter: A pointer to the initialized LogoPassRouter.
func NewLogoPassRouter(h LogoPassHandler, m Middleware) *LogoPassRouter {
	return &LogoPassRouter{
		h: h,
		m: m,
	}
}

// RegisterRoutes registers the routes for logo-password-related operations.
//
// Routes:
//   - POST /api/logo-pass/ - Requires authentication. Calls the Create handler.
//   - GET /api/logo-pass/user/{userID} - Requires authentication. Calls the GetAll handler.
//   - PUT /api/logo-pass/{logoPassID} - Requires authentication. Calls the Update handler.
//
// Parameters:
//   - r chi.Router: The router where the routes will be registered.
func (c *LogoPassRouter) RegisterRoutes(r chi.Router) {
	r.Route("/api/logo-pass", func(r chi.Router) {
		r.With(c.m.Auth).Post("/", c.h.Create)             // Create a new logo-password entry
		r.With(c.m.Auth).Get("/user/{userID}", c.h.GetAll) // Get all logo-passwords for a user
		r.With(c.m.Auth).Put("/{logoPassID}", c.h.Update)  // Update an existing logo-password entry
	})
}
