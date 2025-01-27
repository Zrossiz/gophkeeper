package postgres

import (
	"database/sql"
	"fmt"

	"github.com/Zrossiz/gophkeeper/internal/dto"
	"github.com/Zrossiz/gophkeeper/internal/entities"
)

type LogoPassStorage struct {
	db *sql.DB
}

func NewLogoPassStorage(db *sql.DB) *LogoPassStorage {
	return &LogoPassStorage{db: db}
}

func (l *LogoPassStorage) CreateLogoPass(body dto.CreateLogoPassDTO) error {
	query := `INSERT INTO passwords (user_id, app_name, username, password, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, NOW(), NOW())`
	_, err := l.db.Exec(query, body.UserId, body.AppName, body.Username, body.Password)
	if err != nil {
		return fmt.Errorf("failed to create logo pass: %w", err)
	}
	return nil
}

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
