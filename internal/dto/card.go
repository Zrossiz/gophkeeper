package dto

type CreateCardDTO struct {
	UserID   string `json:"user_id"`
	BankName string `json:"bank_name"`
	Num      string `json:"num"`
	CVV      string `json:"cvv"`
	ExpDate  string `json:"exp_date"`
}

type UpdateCardDTO struct {
	Num     string `json:"num"`
	CVV     string `json:"cvv"`
	ExpDate string `json:"exp_date"`
}
