// Package postgres provides the data storage implementation for handling notes in a PostgreSQL database.
package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Zrossiz/gophkeeper/internal/dto"
	"github.com/Zrossiz/gophkeeper/internal/entities"
)

// NotesStorage represents the storage layer for managing notes in the database.
type NotesStorage struct {
	db *sql.DB
}

// NewNotesStorage creates a new instance of NotesStorage.
//
// Parameters:
//   - db *sql.DB: a database connection.
//
// Returns:
//   - *NotesStorage: a pointer to a NotesStorage instance.
func NewNotesStorage(db *sql.DB) *NotesStorage {
	return &NotesStorage{db: db}
}

// Create inserts a new note into the database.
//
// Parameters:
//   - body dto.CreateNoteDTO: data transfer object containing the note details.
//
// Returns:
//   - error: an error if the insertion fails, otherwise nil.
func (n *NotesStorage) Create(ctx context.Context, body dto.CreateNoteDTO) error {
	query := `INSERT INTO notes (user_id, title, text_data) VALUES ($1, $2, $3)`

	_, err := n.db.ExecContext(ctx, query, body.UserID, body.Title, body.TextData)
	if err != nil {
		return fmt.Errorf("failed to create note: %w", err)
	}

	return nil
}

// Update modifies an existing note in the database.
//
// Parameters:
//   - noteID int: the ID of the note to be updated.
//   - body dto.UpdateNoteDTO: data transfer object containing the updated note details.
//
// Returns:
//   - error: an error if the update fails, otherwise nil.
func (n *NotesStorage) Update(ctx context.Context, noteID int, body dto.UpdateNoteDTO) error {
	query := `UPDATE notes SET title = $1, text_data = $2, updated_at = NOW() WHERE id = $3`

	_, err := n.db.ExecContext(ctx, query, body.Title, body.TextData, noteID)
	if err != nil {
		return fmt.Errorf("failed to update note: %w", err)
	}

	return nil
}

// GetAllByUser retrieves all notes associated with a specific user.
//
// Parameters:
//   - userID int: the ID of the user whose notes should be retrieved.
//
// Returns:
//   - []entities.Note: a slice of Note entities.
//   - error: an error if the retrieval fails, otherwise nil.
func (n *NotesStorage) GetAllByUser(ctx context.Context, userID int) ([]entities.Note, error) {
	query := `SELECT id, user_id, title, text_data, created_at, updated_at FROM notes WHERE user_id = $1`

	rows, err := n.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get all notes: %w", err)
	}
	defer rows.Close()

	var notes []entities.Note
	for rows.Next() {
		var note entities.Note
		err := rows.Scan(
			&note.ID,
			&note.UserID,
			&note.Title,
			&note.TextData,
			&note.CreatedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		notes = append(notes, note)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return notes, nil
}
