package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zrossiz/gophkeeper/internal/apperrors"
	"github.com/Zrossiz/gophkeeper/internal/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// Mock UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Registration(registrationDTO dto.UserDTO) (*dto.GeneratedJwt, error) {
	args := m.Called(registrationDTO)
	if jwt, ok := args.Get(0).(*dto.GeneratedJwt); ok {
		return jwt, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserService) Login(loginDTO dto.UserDTO) (*dto.GeneratedJwt, error) {
	args := m.Called(loginDTO)
	if jwt, ok := args.Get(0).(*dto.GeneratedJwt); ok {
		return jwt, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestUserHandler_Registration_Success(t *testing.T) {
	mockService := new(MockUserService)
	logger := zap.NewNop()
	handler := NewUserHandler(mockService, logger)

	userData := dto.UserDTO{
		Username: "testuser",
		Password: "password123",
	}

	mockJWT := &dto.GeneratedJwt{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		Hash:         "hash-value",
	}

	mockService.On("Registration", userData).Return(mockJWT, nil)

	body, _ := json.Marshal(userData)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Registration(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"hash":"hash-value"`)
}

func TestUserHandler_Registration_EmptyUsername(t *testing.T) {
	mockService := new(MockUserService)
	logger := zap.NewNop()
	handler := NewUserHandler(mockService, logger)

	userData := dto.UserDTO{
		Username: "",
		Password: "password123",
	}

	body, _ := json.Marshal(userData)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Registration(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "login can not be empty")
}

func TestUserHandler_Registration_EmptyPassword(t *testing.T) {
	mockService := new(MockUserService)
	logger := zap.NewNop()
	handler := NewUserHandler(mockService, logger)

	userData := dto.UserDTO{
		Username: "testuser",
		Password: "",
	}

	body, _ := json.Marshal(userData)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Registration(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "password can not be empty")
}

func TestUserHandler_Registration_UserExists(t *testing.T) {
	mockService := new(MockUserService)
	logger := zap.NewNop()
	handler := NewUserHandler(mockService, logger)

	userData := dto.UserDTO{
		Username: "existinguser",
		Password: "password123",
	}

	mockService.On("Registration", userData).Return(nil, apperrors.ErrUserAlreadyExists)

	body, _ := json.Marshal(userData)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Registration(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)
	assert.Contains(t, rec.Body.String(), apperrors.ErrUserAlreadyExists.Error())
}

func TestUserHandler_Registration_DBError(t *testing.T) {
	mockService := new(MockUserService)
	logger := zap.NewNop()
	handler := NewUserHandler(mockService, logger)

	userData := dto.UserDTO{
		Username: "testuser",
		Password: "password123",
	}

	mockService.On("Registration", userData).Return(nil, apperrors.ErrDBQuery)

	body, _ := json.Marshal(userData)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Registration(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "internal server error")
}

func TestUserHandler_Login_Success(t *testing.T) {
	mockService := new(MockUserService)
	logger := zap.NewNop()
	handler := NewUserHandler(mockService, logger)

	userData := dto.UserDTO{
		Username: "testuser",
		Password: "password123",
	}

	mockJWT := &dto.GeneratedJwt{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		Hash:         "hash-value",
	}

	mockService.On("Login", userData).Return(mockJWT, nil)

	body, _ := json.Marshal(userData)
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestUserHandler_Login_InvalidPassword(t *testing.T) {
	mockService := new(MockUserService)
	logger := zap.NewNop()
	handler := NewUserHandler(mockService, logger)

	userData := dto.UserDTO{
		Username: "testuser",
		Password: "wrongpassword",
	}

	mockService.On("Login", userData).Return(nil, apperrors.ErrInvalidPassword)

	body, _ := json.Marshal(userData)
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "unauthorized")
}

func TestUserHandler_Login_UserNotFound(t *testing.T) {
	mockService := new(MockUserService)
	logger := zap.NewNop()
	handler := NewUserHandler(mockService, logger)

	userData := dto.UserDTO{
		Username: "nonexistentuser",
		Password: "password123",
	}

	mockService.On("Login", userData).Return(nil, apperrors.ErrUserAlreadyExists)

	body, _ := json.Marshal(userData)
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "user not found")
}
