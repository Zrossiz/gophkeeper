package entities

import "time"

type Card struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	BankName  string    `json:"bank_name"`
	Number    string    `json:"num"`
	CVV       string    `json:"cvv"`
	ExpDate   string    `json:"exp_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
