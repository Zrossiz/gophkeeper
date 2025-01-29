package postgres

import (
	"database/sql"
	"fmt"

	"github.com/Zrossiz/gophkeeper/internal/apperrors"
	"github.com/Zrossiz/gophkeeper/internal/dto"
	"github.com/Zrossiz/gophkeeper/internal/entities"
)

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{}
}

func (u *UserStorage) Create(body dto.UserDTO) error {
	query := `INSERT INTO users (username, password) VALUES ($1, $2)`
	_, err := u.db.Exec(query, body.Username, body.Password)
	if err != nil {
		return fmt.Errorf("create user error: %v", err)
	}

	return nil
}

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
