package dto

type CreateBinaryDTO struct {
	UserId string `json:"user_id"`
	Title  string `json:"title"`
	Data   string `json:"data"`
}

type UpdateBinaryDTO struct {
	Data string `json:"data"`
}
