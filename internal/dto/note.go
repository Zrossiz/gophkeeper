package dto

type CreateNoteDTO struct {
	UserID   int    `json:"user_id"`
	Title    string `json:"title"`
	TextData string `json:"text_data"`
	Key      string
}

type UpdateNoteDTO struct {
	Title    string `json:"title"`
	TextData string `json:"text_data"`
	Key      string
}
