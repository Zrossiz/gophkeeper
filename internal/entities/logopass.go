package entities

import "time"

type LogoPassword struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	AppName   string    `json:"app_name"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
