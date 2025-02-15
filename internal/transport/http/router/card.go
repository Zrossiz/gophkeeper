// Package router defines the HTTP routing structure for handling card-related requests.
package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// CardRouter provides route registration for card-related HTTP handlers.
type CardRouter struct {
	h CardHandler // Handler for card-related operations.
	m Middleware  // Middleware for authentication and request processing.
}

// CardHandler defines the interface for handling card data requests.
type CardHandler interface {
	// GetAll retrieves all stored card information for a specific user.
	GetAll(rw http.ResponseWriter, r *http.Request)

	// Update modifies an existing card entry.
	Update(rw http.ResponseWriter, r *http.Request)

	// Create adds a new card entry to the storage.
	Create(rw http.ResponseWriter, r *http.Request)
}

// NewCardRouter initializes a new CardRouter instance.
//
// Parameters:
//   - h CardHandler: The handler for card operations.
//   - m Middleware: Middleware for handling authentication and authorization.
//
// Returns:
//   - *CardRouter: A pointer to the initialized CardRouter.
func NewCardRouter(h CardHandler, m Middleware) *CardRouter {
	return &CardRouter{
		h: h,
		m: m,
	}
}

// RegisterRoutes registers the routes for card-related operations.
//
// Routes:
//   - POST /api/card/ - Requires authentication. Calls the Create handler.
//   - GET /api/card/user/{userID} - Requires authentication. Calls the GetAll handler.
//   - PUT /api/card/{cardID} - Requires authentication. Calls the Update handler.
//
// Parameters:
//   - r chi.Router: The router where the routes will be registered.
func (c *CardRouter) RegisterRoutes(r chi.Router) {
	r.Route("/api/card", func(r chi.Router) {
		r.With(c.m.Auth).Post("/", c.h.Create)             // Create a new card entry
		r.With(c.m.Auth).Get("/user/{userID}", c.h.GetAll) // Get all cards for a user
		r.With(c.m.Auth).Put("/{cardID}", c.h.Update)      // Update an existing card entry
	})
}
