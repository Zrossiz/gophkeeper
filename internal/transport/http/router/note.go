// Package router defines the HTTP routing structure for handling note-related requests.
package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// NoteRouter provides route registration for note-related HTTP handlers.
type NoteRouter struct {
	h NoteHandler // Handler for note operations.
	m Middleware  // Middleware for authentication and request processing.
}

// NoteHandler defines the interface for handling note data requests.
type NoteHandler interface {
	// GetAll retrieves all stored notes for a specific user.
	GetAll(rw http.ResponseWriter, r *http.Request)

	// Update modifies an existing note entry.
	Update(rw http.ResponseWriter, r *http.Request)

	// Create adds a new note entry to the storage.
	Create(rw http.ResponseWriter, r *http.Request)
}

// NewNoteRouter initializes a new NoteRouter instance.
//
// Parameters:
//   - h NoteHandler: The handler for note operations.
//   - m Middleware: Middleware for handling authentication and authorization.
//
// Returns:
//   - *NoteRouter: A pointer to the initialized NoteRouter.
func NewNoteRouter(h NoteHandler, m Middleware) *NoteRouter {
	return &NoteRouter{
		h: h,
		m: m,
	}
}

// RegisterRoutes registers the routes for note-related operations.
//
// Routes:
//   - POST /api/note/ - Requires authentication. Calls the Create handler.
//   - GET /api/note/user/{userID} - Requires authentication. Calls the GetAll handler.
//   - PUT /api/note/{noteID} - Requires authentication. Calls the Update handler.
//
// Parameters:
//   - r chi.Router: The router where the routes will be registered.
func (n *NoteRouter) RegisterRoutes(r chi.Router) {
	r.Route("/api/note", func(r chi.Router) {
		r.With(n.m.Auth).Post("/", n.h.Create)             // Create a new note
		r.With(n.m.Auth).Get("/user/{userID}", n.h.GetAll) // Get all notes for a user
		r.With(n.m.Auth).Put("/{noteID}", n.h.Update)      // Update an existing note
	})
}
