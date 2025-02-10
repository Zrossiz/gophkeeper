package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zrossiz/gophkeeper/internal/dto"
	"github.com/Zrossiz/gophkeeper/internal/entities"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// Mock сервиса
type MockLogoPassService struct {
	mock.Mock
}

func (m *MockLogoPassService) Create(ctx context.Context, body dto.CreateLogoPassDTO) error {
	args := m.Called(body)
	return args.Error(0)
}

func (m *MockLogoPassService) Update(ctx context.Context, userID int64, body dto.UpdateLogoPassDTO) error {
	args := m.Called(userID, body)
	return args.Error(0)
}

func (m *MockLogoPassService) GetAll(ctx context.Context, userID int64, key string) ([]entities.LogoPassword, error) {
	args := m.Called(userID, key)
	return args.Get(0).([]entities.LogoPassword), args.Error(1)
}

func setupTestHandler() (*LogoPassHandler, *MockLogoPassService) {
	mockService := new(MockLogoPassService)
	logger := zap.NewNop()
	handler := NewLogoPassHandler(mockService, logger)
	return handler, mockService
}

func TestCreate_Success(t *testing.T) {
	handler, mockService := setupTestHandler()
	mockService.On("Create", mock.Anything).Return(nil)

	body := dto.CreateLogoPassDTO{
		Username: "testuser",
		Password: "testpass",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/logo-pass", bytes.NewBuffer(bodyBytes))
	req.AddCookie(&http.Cookie{Name: "key", Value: "testkey"})
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	mockService.AssertExpectations(t)
}

func TestCreate_MissingCookie(t *testing.T) {
	handler, _ := setupTestHandler()

	req := httptest.NewRequest("POST", "/logo-pass", nil)
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "key not found")
}

func TestCreate_InvalidJSON(t *testing.T) {
	handler, _ := setupTestHandler()

	req := httptest.NewRequest("POST", "/logo-pass", bytes.NewBuffer([]byte("{invalid")))
	req.AddCookie(&http.Cookie{Name: "key", Value: "testkey"})
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreate_ServiceError(t *testing.T) {
	handler, mockService := setupTestHandler()
	mockService.On("Create", mock.Anything).Return(errors.New("internal error"))

	body := dto.CreateLogoPassDTO{Username: "testuser", Password: "testpass"}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/logo-pass", bytes.NewBuffer(bodyBytes))
	req.AddCookie(&http.Cookie{Name: "key", Value: "testkey"})
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockService.AssertExpectations(t)
}

func TestUpdate_Success(t *testing.T) {
	handler, mockService := setupTestHandler()
	mockService.On("Update", int64(1), mock.Anything).Return(nil)

	body := dto.UpdateLogoPassDTO{Password: "newpass"}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("PUT", "/logo-pass/1", bytes.NewBuffer(bodyBytes))
	req.AddCookie(&http.Cookie{Name: "key", Value: "testkey"})

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("logoPassID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rec := httptest.NewRecorder()
	handler.Update(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestUpdate_MissingCookie(t *testing.T) {
	handler, _ := setupTestHandler()

	req := httptest.NewRequest("PUT", "/logo-pass/1", nil)
	rec := httptest.NewRecorder()

	handler.Update(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdate_InvalidID(t *testing.T) {
	handler, _ := setupTestHandler()

	req := httptest.NewRequest("PUT", "/logo-pass/abc", nil)
	rec := httptest.NewRecorder()

	handler.Update(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetAll_MissingCookie(t *testing.T) {
	handler, _ := setupTestHandler()

	req := httptest.NewRequest("GET", "/logo-pass/user/1", nil)
	rec := httptest.NewRecorder()

	handler.GetAll(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetAll_InvalidUserID(t *testing.T) {
	handler, _ := setupTestHandler()

	req := httptest.NewRequest("GET", "/logo-pass/user/abc", nil)
	rec := httptest.NewRecorder()

	handler.GetAll(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
