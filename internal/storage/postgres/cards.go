// Package postgres provides database storage implementations for various entities.
package postgres

import (
	"database/sql"

	"github.com/Zrossiz/gophkeeper/internal/dto"
	"github.com/Zrossiz/gophkeeper/internal/entities"
)

// CardStorage provides methods for managing card-related data in the PostgreSQL database.
type CardStorage struct {
	db *sql.DB // SQL database connection.
}

// NewCardStorage initializes and returns a new CardStorage instance.
//
// Parameters:
//   - db: An active SQL database connection.
//
// Returns:
//   - A pointer to an initialized CardStorage instance.
func NewCardStorage(db *sql.DB) *CardStorage {
	return &CardStorage{db: db}
}

// CreateCard inserts a new card record into the database.
//
// Parameters:
//   - body: A CreateCardDTO struct containing card details such as bank name, number, CVV, expiration date, and cardholder name.
//
// Returns:
//   - An error if the operation fails.
func (c *CardStorage) CreateCard(body dto.CreateCardDTO) error {
	query := `
		INSERT INTO cards (user_id, bank_name, num, cvv, exp_date, card_holder_name) 
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := c.db.Exec(
		query,
		body.UserID,
		body.BankName,
		body.Num,
		body.CVV,
		body.ExpDate,
		body.CardHolderName,
	)
	if err != nil {
		return err
	}

	return nil
}

// GetAllCardsByUserId retrieves all stored cards for a given user.
//
// Parameters:
//   - userID: The unique identifier of the user.
//
// Returns:
//   - A slice of Card entities containing the user's stored cards.
//   - An error if the retrieval fails.
func (c *CardStorage) GetAllCardsByUserId(userID int64) ([]entities.Card, error) {
	query := `SELECT id, user_id, bank_name, num, cvv, exp_date, card_holder_name, created_at, updated_at 
              FROM cards WHERE user_id = $1`

	rows, err := c.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []entities.Card
	for rows.Next() {
		var card entities.Card
		err := rows.Scan(
			&card.ID,
			&card.UserID,
			&card.BankName,
			&card.Number,
			&card.CVV,
			&card.ExpDate,
			&card.CardHolderName,
			&card.CreatedAt,
			&card.UpdatedAt,
		)
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

// UpdateCard modifies an existing card record in the database.
//
// Parameters:
//   - cardID: The unique identifier of the card to be updated.
//   - body: An UpdateCardDTO struct containing the updated card details.
//
// Returns:
//   - An error if the update operation fails.
func (c *CardStorage) UpdateCard(cardID int64, body dto.UpdateCardDTO) error {
	query := `UPDATE cards 
              SET num = $1, cvv = $2, exp_date = $3, card_holder_name = $4, updated_at = NOW() 
              WHERE id = $5`

	_, err := c.db.Exec(query, body.Num, body.CVV, body.ExpDate, body.CardHolderName, cardID)
	if err != nil {
		return err
	}

	return nil
}
