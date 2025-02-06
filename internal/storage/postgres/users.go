// Package postgres provides database storage implementations for handling user-related data.
package postgres

import (
	"database/sql"
	"fmt"

	"github.com/Zrossiz/gophkeeper/internal/apperrors"
	"github.com/Zrossiz/gophkeeper/internal/dto"
	"github.com/Zrossiz/gophkeeper/internal/entities"
)

// UserStorage handles database operations related to users.
type UserStorage struct {
	db *sql.DB // Database connection instance.
}

// NewUserStorage initializes a new UserStorage instance.
//
// Parameters:
//   - db *sql.DB: The database connection.
//
// Returns:
//   - *UserStorage: A pointer to the initialized UserStorage structure.
func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

// Create inserts a new user record into the database.
//
// Parameters:
//   - body dto.UserDTO: The user data transfer object containing the username and password.
//
// Returns:
//   - error: An error if the operation fails, otherwise nil.
func (u *UserStorage) Create(body dto.UserDTO) error {
	query := `INSERT INTO users (username, password) VALUES ($1, $2)`
	_, err := u.db.Exec(query, body.Username, body.Password)
	if err != nil {
		return fmt.Errorf("create user error: %v", err)
	}

	return nil
}

// GetUserByUsername retrieves a user record by their username.
//
// Parameters:
//   - username string: The username of the user to retrieve.
//
// Returns:
//   - *entities.User: A pointer to the retrieved user entity if found.
//   - error: Returns an error if the user is not found or if a query error occurs.
func (u *UserStorage) GetUserByUsername(username string) (*entities.User, error) {
	query := `SELECT id, username, password FROM users WHERE username = $1`
	row := u.db.QueryRow(query, username)
	var user entities.User
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
