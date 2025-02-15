// Package router defines the HTTP routing structure for handling binary data requests.
package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// BinaryRouter provides route registration for binary data-related HTTP handlers.
type BinaryRouter struct {
	h BinaryHandler // Handler for binary-related operations.
	m Middleware    // Middleware for authentication and request processing.
}

// BinaryHandler defines the interface for handling binary data requests.
type BinaryHandler interface {
	// GetAll retrieves all binary data for a specific user.
	GetAll(rw http.ResponseWriter, r *http.Request)

	// Create adds new binary data to storage.
	Create(rw http.ResponseWriter, r *http.Request)
}

// NewBinaryRouter initializes a new BinaryRouter instance.
//
// Parameters:
//   - h BinaryHandler: The handler for binary data operations.
//   - m Middleware: Middleware for handling authentication and authorization.
//
// Returns:
//   - *BinaryRouter: A pointer to the initialized BinaryRouter.
func NewBinaryRouter(h BinaryHandler, m Middleware) *BinaryRouter {
	return &BinaryRouter{
		h: h,
		m: m,
	}
}

// RegisterRoutes registers the routes for binary data operations.
//
// Routes:
//   - POST /api/binary/ - Requires authentication. Calls the Create handler.
//   - GET /api/binary/user/{userID} - Requires authentication. Calls the GetAll handler.
//
// Parameters:
//   - r chi.Router: The router where the routes will be registered.
func (b *BinaryRouter) RegisterRoutes(r chi.Router) {
	r.Route("/api/binary", func(r chi.Router) {
		r.With(b.m.Auth).Post("/", b.h.Create)             // Create binary data
		r.With(b.m.Auth).Get("/user/{userID}", b.h.GetAll) // Get all binary data for a user
	})
}
