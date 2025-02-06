// Package postgres provides database storage implementations for various entities.
package postgres

import (
	"database/sql"
	"fmt"

	"github.com/Zrossiz/gophkeeper/internal/dto"
	"github.com/Zrossiz/gophkeeper/internal/entities"
)

// LogoPassStorage handles database operations related to stored application passwords.
type LogoPassStorage struct {
	db *sql.DB // SQL database connection.
}

// NewLogoPassStorage initializes and returns a new LogoPassStorage instance.
//
// Parameters:
//   - db: An active SQL database connection.
//
// Returns:
//   - A pointer to an initialized LogoPassStorage instance.
func NewLogoPassStorage(db *sql.DB) *LogoPassStorage {
	return &LogoPassStorage{db: db}
}

// CreateLogoPass inserts a new application password record into the database.
//
// Parameters:
//   - body: A CreateLogoPassDTO struct containing user ID, application name, username, and password.
//
// Returns:
//   - An error if the operation fails.
func (l *LogoPassStorage) CreateLogoPass(body dto.CreateLogoPassDTO) error {
	query := `INSERT INTO passwords (user_id, app_name, username, password, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, NOW(), NOW())`
	_, err := l.db.Exec(query, body.UserId, body.AppName, body.Username, body.Password)
	if err != nil {
		return fmt.Errorf("failed to create logo pass: %w", err)
	}
	return nil
}

// GetAllByUser retrieves all stored application passwords for a given user.
//
// Parameters:
//   - userID: The unique identifier of the user.
//
// Returns:
//   - A slice of LogoPassword entities containing the user's stored credentials.
//   - An error if the retrieval fails.
func (l *LogoPassStorage) GetAllByUser(userID int64) ([]entities.LogoPassword, error) {
	query := `SELECT id, user_id, app_name, username, password, created_at, updated_at 
              FROM passwords WHERE user_id = $1`

	rows, err := l.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get logo passes: %w", err)
	}
	defer rows.Close()

	var logoPasswords []entities.LogoPassword
	for rows.Next() {
		var lp entities.LogoPassword
		err := rows.Scan(&lp.ID, &lp.UserID, &lp.AppName, &lp.Username, &lp.Password, &lp.CreatedAt, &lp.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		logoPasswords = append(logoPasswords, lp)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return logoPasswords, nil
}

// UpdateLogoPass modifies an existing application password record in the database.
//
// Parameters:
//   - id: The unique identifier of the password record to be updated.
//   - body: An UpdateLogoPassDTO struct containing the updated username and password.
//
// Returns:
//   - An error if the update operation fails.
func (l *LogoPassStorage) UpdateLogoPass(id int64, body dto.UpdateLogoPassDTO) error {
	query := `UPDATE passwords 
              SET username = $1, password = $2, updated_at = NOW() 
              WHERE id = $3`
	_, err := l.db.Exec(query, body.Username, body.Password, id)
	if err != nil {
		return fmt.Errorf("failed to update logo pass: %w", err)
	}
	return nil
}
