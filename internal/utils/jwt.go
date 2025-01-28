package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	UserID   int64  `json:"userID"`
	Username string `json:"userName"`
	jwt.RegisteredClaims
}

type GenerateJWTProps struct {
	Secret   []byte
	Exprires time.Time
	UserID   int64
	Username string
}

func GenerateJWT(props GenerateJWTProps) (string, error) {
	claims := &CustomClaims{
		UserID:   props.UserID,
		Username: props.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(props.Exprires),
			Issuer:    "exampleIssuer",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(props.Secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
