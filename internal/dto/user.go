package dto

type UserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type GeneratedJwt struct {
	AccessToken  string
	RefreshToken string
	Hash         string `json:"hash"`
}
