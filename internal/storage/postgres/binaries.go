package postgres

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Zrossiz/gophkeeper/internal/dto"
	"github.com/Zrossiz/gophkeeper/internal/entities"
)

type BinaryStorage struct {
	db *sql.DB
}

func NewBinaryStorage(db *sql.DB) *BinaryStorage {
	return &BinaryStorage{db: db}
}

func (b *BinaryStorage) Create(body dto.CreateBinaryDTO) error {
	query := `
		INSERT INTO binary_data (user_id, title, binary_data, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := b.db.Exec(query, body.UserId, body.Title, body.Data, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (b *BinaryStorage) Update(id int64, body dto.UpdateBinaryDTO) error {
	query := `
		UPDATE binary_data
		SET binary_data = $1, updated_at = $2
		WHERE id = $3
	`
	result, err := b.db.Exec(query, body.Data, time.Now(), id)
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
