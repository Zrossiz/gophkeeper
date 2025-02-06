// Package postgres provides the database storage implementation for various entities,
// including users, cards, login credentials, binary data, and notes.
package postgres

import (
	"database/sql"
)

// Storage aggregates all storage components used for handling different types of data.
type Storage struct {
	Binary   BinaryStorage   // Handles storage operations for binary data.
	Card     CardStorage     // Manages storage operations for card-related data.
	LogoPass LogoPassStorage // Stores login credentials (username, password).
	User     UserStorage     // Manages user-related storage operations.
	Note     NotesStorage    // Handles note storage operations.
}

// New initializes a new Storage instance with the provided database connection.
//
// Parameters:
//   - conn *sql.DB: The active database connection.
//
// Returns:
//   - *Storage: A pointer to the initialized Storage structure.
func New(conn *sql.DB) *Storage {
	return &Storage{
		User:     *NewUserStorage(conn),
		Card:     *NewCardStorage(conn),
		LogoPass: *NewLogoPassStorage(conn),
		Binary:   *NewBinaryStorage(conn),
		Note:     *NewNotesStorage(conn),
	}
}

// Connect establishes a connection to the PostgreSQL database.
//
// Parameters:
//   - DBURI string: The database connection string.
//
// Returns:
//   - *sql.DB: A pointer to the established database connection.
//   - error: An error if the connection fails, otherwise nil.
func Connect(DBURI string) (*sql.DB, error) {
	db, err := sql.Open("postgres", DBURI)
	if err != nil {
		return nil, err
	}

	return db, nil
}
