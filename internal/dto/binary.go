package dto

type CreateBinaryDTO struct {
	UserID int    `json:"user_id"`
	Title  string `json:"title"`
	Data   []byte `json:"data"`
	Key    string
}

type SetStorageBinaryDTO struct {
	UserID int    `json:"user_id"`
	Title  string `json:"title"`
	Data   []byte `json:"data"`
}

type UpdateBinaryDTO struct {
	Title string `json:"title"`
	Data  []byte `json:"data"`
}
