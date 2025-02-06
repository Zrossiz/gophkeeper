// Package postgres provides database storage implementations for various entities.
package postgres

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Zrossiz/gophkeeper/internal/dto"
	"github.com/Zrossiz/gophkeeper/internal/entities"
)

// BinaryStorage provides methods for managing binary data in the PostgreSQL database.
type BinaryStorage struct {
	db *sql.DB // SQL database connection.
}

// NewBinaryStorage initializes and returns a new BinaryStorage instance.
//
// Parameters:
//   - db: An active SQL database connection.
//
// Returns:
//   - A pointer to an initialized BinaryStorage instance.
func NewBinaryStorage(db *sql.DB) *BinaryStorage {
	return &BinaryStorage{db: db}
}

// Create inserts a new binary data record into the database.
//
// Parameters:
//   - body: A SetStorageBinaryDTO struct containing user ID, title, and binary data.
//
// Returns:
//   - An error if the operation fails.
func (b *BinaryStorage) Create(body dto.SetStorageBinaryDTO) error {
	query := `
		INSERT INTO binary_data (user_id, title, binary_data, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := b.db.Exec(query, body.UserID, body.Title, body.Data, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}

// Update modifies an existing binary data record in the database.
//
// Parameters:
//   - body: A SetStorageBinaryDTO struct containing the updated binary data.
//
// Returns:
//   - An error if the update operation fails or the record is not found.
func (b *BinaryStorage) Update(body dto.SetStorageBinaryDTO) error {
	query := `
		UPDATE binary_data
		SET binary_data = $1, updated_at = $2
		WHERE id = $3
	`
	result, err := b.db.Exec(query, body.Data, time.Now(), body.UserID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no rows affected, record not found")
	}

	return nil
}

// GetAllByUser retrieves all binary data records for a given user.
//
// Parameters:
//   - userID: The unique identifier of the user.
//
// Returns:
//   - A slice of BinaryData entities containing the user's stored binary data.
//   - An error if the retrieval fails.
func (b *BinaryStorage) GetAllByUser(userID int64) ([]entities.BinaryData, error) {
	query := `
		SELECT id, user_id, title, binary_data, created_at, updated_at
		FROM binary_data
		WHERE user_id = $1
	`
	rows, err := b.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var binaryDataList []entities.BinaryData
	for rows.Next() {
		var binaryData entities.BinaryData
		err := rows.Scan(
			&binaryData.ID,
			&binaryData.UserID,
			&binaryData.Title,
			&binaryData.Data,
			&binaryData.CreatedAt,
			&binaryData.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		binaryDataList = append(binaryDataList, binaryData)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return binaryDataList, nil
}
