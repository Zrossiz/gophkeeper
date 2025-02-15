package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Zrossiz/gophkeeper/internal/config"
	"github.com/Zrossiz/gophkeeper/internal/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func generateValidToken(secret string, userID int64, username string, duration time.Duration) (string, error) {
	return utils.GenerateJWT(utils.GenerateJWTProps{
		Secret:   []byte(secret),
		Exprires: time.Now().Add(duration),
		UserID:   userID,
		Username: username,
	})
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	cfg := config.Config{AccessSecret: "testsecret"}
	logger := zap.NewNop()
	middleware := New(cfg, logger)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(UserIDContextKey).(int64)
		assert.True(t, ok, "UserID not found in context")
		assert.Equal(t, int64(1), userID)

		username, ok := r.Context().Value(UserNameContextKey).(string)
		assert.True(t, ok, "Username not found in context")
		assert.Equal(t, "testuser", username)

		w.WriteHeader(http.StatusOK)
	})

	validToken, err := generateValidToken(cfg.AccessSecret, 1, "testuser", 1*time.Hour)
	assert.NoError(t, err)

	req := httptest.NewRequest("GET", "http://localhost/", nil)
	req.AddCookie(&http.Cookie{Name: "accesstoken", Value: validToken})
	rec := httptest.NewRecorder()

	middleware.Auth(testHandler).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthMiddleware_NoToken(t *testing.T) {
	cfg := config.Config{AccessSecret: "testsecret"}
	logger := zap.NewNop()
	middleware := New(cfg, logger)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://localhost/", nil)

	middleware.Auth(http.NotFoundHandler()).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "unauthorized")
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	cfg := config.Config{AccessSecret: "testsecret"}
	logger := zap.NewNop()
	middleware := New(cfg, logger)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://localhost/", nil)
	req.AddCookie(&http.Cookie{Name: "accesstoken", Value: "invalid.token.here"})

	middleware.Auth(http.NotFoundHandler()).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "unauthorized")
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	cfg := config.Config{AccessSecret: "testsecret"}
	logger := zap.NewNop()
	middleware := New(cfg, logger)

	expiredToken, err := generateValidToken(cfg.AccessSecret, 1, "testuser", -1*time.Hour)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://localhost/", nil)
	req.AddCookie(&http.Cookie{Name: "accesstoken", Value: expiredToken})

	middleware.Auth(http.NotFoundHandler()).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "unauthorized: invalid token")
}
