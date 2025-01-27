package dto

type CreateLogoPassDTO struct {
	UserId   string `json:"user_id"`
	AppName  string `json:"app_name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateLogoPassDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
