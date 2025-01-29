package postgres

import (
	"database/sql"

	"github.com/Zrossiz/gophkeeper/internal/dto"
	"github.com/Zrossiz/gophkeeper/internal/entities"
)

type CardStorage struct {
	db *sql.DB
}

func NewCardStorage(db *sql.DB) *CardStorage {
	return &CardStorage{db: db}
}

func (c *CardStorage) CreateCard(body dto.CreateCardDTO) error {
	query := `INSERT INTO cards (user_id, bankName, num, cvv, exp_date) VALUES ($1, $2, $3, $4, $5)`

	_, err := c.db.Exec(query, body.UserID, body.BankName, body.Num, body.CVV, body.ExpDate)
	if err != nil {
		return err
	}

	return nil
}

func (c *CardStorage) GetAllCardsByUserId(userID int64) ([]entities.Card, error) {
	query := `SELECT id, user_id, bankName, num, cvv, exp_date, created_at, updated_at 
              FROM cards WHERE user_id = $1`

	rows, err := c.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []entities.Card
	for rows.Next() {
		var card entities.Card
		err := rows.Scan(&card.ID, &card.UserID, &card.BankName, &card.Number, &card.CVV, &card.ExpDate, &card.CreatedAt, &card.UpdatedAt)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cards, nil
}

func (c *CardStorage) UpdateCard(cardID int64, body dto.UpdateCardDTO) error {
	query := `UPDATE cards 
              SET num = $1, cvv = $2, exp_date = $3, updated_at = NOW() 
              WHERE id = $4`

	_, err := c.db.Exec(query, body.Num, body.CVV, body.ExpDate, cardID)
	if err != nil {
		return err
	}

	return nil
}
