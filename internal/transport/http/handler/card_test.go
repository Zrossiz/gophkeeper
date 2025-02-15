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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockCardService struct {
	mock.Mock
}

func (m *MockCardService) Create(ctx context.Context, body dto.CreateCardDTO) error {
	args := m.Called(body)
	return args.Error(0)
}

func (m *MockCardService) Update(ctx context.Context, cardID int64, body dto.UpdateCardDTO) error {
	args := m.Called(cardID, body)
	return args.Error(0)
}

func (m *MockCardService) GetAll(ctx context.Context, userID int64, key string) ([]entities.Card, error) {
	args := m.Called(userID, key)
	return args.Get(0).([]entities.Card), args.Error(1)
}

func setupCardTestHandler() (*CardHandler, *MockCardService) {
	mockService := new(MockCardService)
	logger := zap.NewNop()
	handler := NewCardHandler(mockService, logger)
	return handler, mockService
}

func TestCardCreate_Success(t *testing.T) {
	handler, mockService := setupCardTestHandler()

	// Ожидаем, что сервис вернет успех
	mockService.On("Create", mock.AnythingOfType("dto.CreateCardDTO")).Return(nil)

	body := dto.CreateCardDTO{Num: "1234", ExpDate: "12/25", CVV: "123"}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/card", bytes.NewBuffer(bodyBytes))
	req.AddCookie(&http.Cookie{Name: "key", Value: "testkey"})
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	mockService.AssertExpectations(t)
}

func TestCardCreate_InvalidRequestBody(t *testing.T) {
	handler, _ := setupTestHandler()

	req := httptest.NewRequest("POST", "/card", bytes.NewBufferString("{invalid json}"))
	req.AddCookie(&http.Cookie{Name: "key", Value: "testkey"})
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCardUpdate_InvalidRequestBody(t *testing.T) {
	handler, _ := setupTestHandler()

	// Некорректный JSON
	req := httptest.NewRequest("PUT", "/card/1", bytes.NewBufferString("{invalid json}"))
	req.AddCookie(&http.Cookie{Name: "key", Value: "testkey"})
	rec := httptest.NewRecorder()

	handler.Update(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCardUpdate_InvalidUserID(t *testing.T) {
	handler, mockService := setupTestHandler()

	// Мокируем ошибку на уровне сервиса
	mockService.On("Update", mock.Anything, mock.Anything).Return(nil)

	body := dto.UpdateCardDTO{Num: "5678", ExpDate: "12/24", CVV: "456"}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("PUT", "/card/invalidID", bytes.NewBuffer(bodyBytes))
	req.AddCookie(&http.Cookie{Name: "key", Value: "testkey"})
	rec := httptest.NewRecorder()

	handler.Update(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCardGetAll_InvalidUserID(t *testing.T) {
	handler, mockService := setupTestHandler()

	mockService.On("GetAll", mock.Anything, mock.Anything).Return(nil, errors.New("internal error"))

	req := httptest.NewRequest("GET", "/card/user/invalidID", nil)
	req.AddCookie(&http.Cookie{Name: "key", Value: "testkey"})
	rec := httptest.NewRecorder()

	handler.GetAll(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
