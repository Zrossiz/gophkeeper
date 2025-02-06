// Package utils provides utility functions for handling JWT authentication.
package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// CustomClaims represents the custom claims embedded in a JWT token.
type CustomClaims struct {
	UserID   int64  `json:"userID"`   // User's unique identifier.
	Username string `json:"userName"` // User's username.
	jwt.RegisteredClaims
}

// GenerateJWTProps defines the properties required to generate a JWT token.
type GenerateJWTProps struct {
	Secret   []byte    // Secret key used to sign the JWT.
	Exprires time.Time // Expiration time of the token.
	UserID   int64     // User's unique identifier.
	Username string    // User's username.
}

// GenerateJWT generates a new JWT token using the provided properties.
//
// Parameters:
//   - props GenerateJWTProps: The properties required to generate the token.
//
// Returns:
//   - string: The generated JWT token as a string.
//   - error: An error if token generation fails.
func GenerateJWT(props GenerateJWTProps) (string, error) {
	if len(props.Secret) == 0 {
		return "", fmt.Errorf("secret key cannot be empty")
	}

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
