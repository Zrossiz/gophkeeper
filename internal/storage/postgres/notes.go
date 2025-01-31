package postgres

import (
	"database/sql"
	"fmt"

	"github.com/Zrossiz/gophkeeper/internal/dto"
	"github.com/Zrossiz/gophkeeper/internal/entities"
)

type NotesStorage struct {
	db *sql.DB
}

func NewNotesStorage(db *sql.DB) *NotesStorage {
	return &NotesStorage{db: db}
}

func (n *NotesStorage) Create(body dto.CreateNoteDTO) error {
	query := `INSERT INTO notes (user_id, title, text_data) VALUES ($1, $2, $3)`

	_, err := n.db.Exec(query, body.UserID, body.Title, body.TextData)
	if err != nil {
		return fmt.Errorf("failed to create note: %w", err)
	}

	return nil
}

func (n *NotesStorage) Update(noteID int, body dto.UpdateNoteDTO) error {
	query := `UPDATE notes SET title = $1, text_data = $2, updated_at = NOW()`

	_, err := n.db.Exec(query, body.Title, body.TextData, noteID)
	if err != nil {
		return fmt.Errorf("failed to update note: %w", err)
	}

	return nil
}

func (n *NotesStorage) GetAllByUser(userID int) ([]entities.Note, error) {
	query := `SELECT id, user_id, title, text_data, created_at, updated_at WHERE user_id = $1`

	rows, err := n.db.Query(query, userID)
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
