package dto

type CreateLogoPassDTO struct {
	UserId   int    `json:"user_id"`
	AppName  string `json:"app_name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Key      string
}

type UpdateLogoPassDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Key      string
}
