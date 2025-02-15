package dto

type CreateCardDTO struct {
	UserID         int    `json:"user_id"`
	BankName       string `json:"bank_name"`
	Num            string `json:"num"`
	CVV            string `json:"cvv"`
	ExpDate        string `json:"exp_date"`
	CardHolderName string `json:"card_holder_name"`
	Key            string
}

type UpdateCardDTO struct {
	Num            string `json:"num"`
	CVV            string `json:"cvv"`
	ExpDate        string `json:"exp_date"`
	CardHolderName string `json:"card_holder_name"`
	Key            string
}
