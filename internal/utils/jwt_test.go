package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestGenerateValidJWT(t *testing.T) {
	secret := []byte("testsecret")
	props := GenerateJWTProps{
		Secret:   secret,
		Exprires: time.Now().Add(time.Hour),
		UserID:   123,
		Username: "testuser",
	}

	tokenString, err := GenerateJWT(props)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	assert.NoError(t, err)
	assert.True(t, token.Valid)

	claims, ok := token.Claims.(*CustomClaims)
	assert.True(t, ok)
	assert.Equal(t, props.UserID, claims.UserID)
	assert.Equal(t, props.Username, claims.Username)
	assert.Equal(t, "exampleIssuer", claims.Issuer)
	assert.WithinDuration(t, props.Exprires, claims.ExpiresAt.Time, time.Second)
}

func TestGenerateJWTWithEmptySecret(t *testing.T) {
	props := GenerateJWTProps{
		Secret:   []byte{},
		Exprires: time.Now().Add(time.Hour),
		UserID:   123,
		Username: "testuser",
	}

	tokenString, err := GenerateJWT(props)
	assert.Error(t, err)
	assert.Empty(t, tokenString)
}

func TestGenerateJWTWithExpiredTime(t *testing.T) {
	secret := []byte("testsecret")
	props := GenerateJWTProps{
		Secret:   secret,
		Exprires: time.Now().Add(-time.Hour),
		UserID:   123,
		Username: "testuser",
	}

	tokenString, err := GenerateJWT(props)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token is expired")

	assert.NotNil(t, token)
	assert.False(t, token.Valid)
}
