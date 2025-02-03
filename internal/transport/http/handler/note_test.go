package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zrossiz/gophkeeper/internal/apperrors"
	"github.com/Zrossiz/gophkeeper/internal/dto"
	"github.com/Zrossiz/gophkeeper/internal/entities"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockNoteService struct {
	mock.Mock
}

func (m *MockNoteService) Create(body dto.CreateNoteDTO) error {
	args := m.Called(body)
	return args.Error(0)
}

func (m *MockNoteService) Update(noteID int, body dto.UpdateNoteDTO) error {
	args := m.Called(noteID, body)
	return args.Error(0)
}

func (m *MockNoteService) GetAll(userID int, key string) ([]entities.Note, error) {
	args := m.Called(userID, key)
	if notes, ok := args.Get(0).([]entities.Note); ok {
		return notes, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestNoteHandler_Create_Success(t *testing.T) {
	mockService := new(MockNoteService)
	logger := zap.NewNop()
	handler := NewNoteHandler(mockService, logger)

	noteData := dto.CreateNoteDTO{
		Title:    "Test Note",
		TextData: "This is a test note",
		Key:      "test-key",
	}

	mockService.On("Create", noteData).Return(nil)

	body, _ := json.Marshal(noteData)
	req := httptest.NewRequest(http.MethodPost, "/note", bytes.NewReader(body))
	req.AddCookie(&http.Cookie{Name: "key", Value: "test-key"})
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestNoteHandler_Create_MissingKeyCookie(t *testing.T) {
	mockService := new(MockNoteService)
	logger := zap.NewNop()
	handler := NewNoteHandler(mockService, logger)

	noteData := dto.CreateNoteDTO{
		Title:    "Test Note",
		TextData: "This is a test note",
	}

	body, _ := json.Marshal(noteData)
	req := httptest.NewRequest(http.MethodPost, "/note", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "key not found")
}

func TestNoteHandler_Create_InvalidJSON(t *testing.T) {
	mockService := new(MockNoteService)
	logger := zap.NewNop()
	handler := NewNoteHandler(mockService, logger)

	invalidJSON := `{"title": "Test Note", "body":}`

	req := httptest.NewRequest(http.MethodPost, "/note", bytes.NewReader([]byte(invalidJSON)))
	req.AddCookie(&http.Cookie{Name: "key", Value: "test-key"})
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), apperrors.ErrInvalidRequestBody)
}

func TestNoteHandler_Create_InternalServerError(t *testing.T) {
	mockService := new(MockNoteService)
	logger := zap.NewNop()
	handler := NewNoteHandler(mockService, logger)

	noteData := dto.CreateNoteDTO{
		Title:    "Test Note",
		TextData: "This is a test note",
		Key:      "test-key",
	}

	mockService.On("Create", noteData).Return(errors.New("db error"))

	body, _ := json.Marshal(noteData)
	req := httptest.NewRequest(http.MethodPost, "/note", bytes.NewReader(body))
	req.AddCookie(&http.Cookie{Name: "key", Value: "test-key"})
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), apperrors.ErrInternalServer)
}

func TestNoteHandler_Update_Success(t *testing.T) {
	mockService := new(MockNoteService)
	logger := zap.NewNop()
	handler := NewNoteHandler(mockService, logger)

	updateData := dto.UpdateNoteDTO{
		Title:    "Updated Title",
		TextData: "Updated Body",
		Key:      "test-key",
	}

	mockService.On("Update", 1, updateData).Return(nil)

	body, _ := json.Marshal(updateData)
	req := httptest.NewRequest(http.MethodPut, "/note/1", bytes.NewReader(body))
	req.AddCookie(&http.Cookie{Name: "key", Value: "test-key"})
	req.Header.Set("Content-Type", "application/json")

	// Chi router нужен для парсинга параметров
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("noteID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rec := httptest.NewRecorder()

	handler.Update(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestNoteHandler_Update_InvalidNoteID(t *testing.T) {
	mockService := new(MockNoteService)
	logger := zap.NewNop()
	handler := NewNoteHandler(mockService, logger)

	req := httptest.NewRequest(http.MethodPut, "/note/abc", nil)
	rec := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("noteID", "abc")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.Update(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid note id")
}

func TestNoteHandler_GetAll_Success(t *testing.T) {
	mockService := new(MockNoteService)
	logger := zap.NewNop()
	handler := NewNoteHandler(mockService, logger)

	mockNotes := []entities.Note{
		{ID: 1, Title: "Note 1", TextData: "Body 1"},
		{ID: 2, Title: "Note 2", TextData: "Body 2"},
	}

	mockService.On("GetAll", 1, "test-key").Return(mockNotes, nil)

	req := httptest.NewRequest(http.MethodGet, "/note/user/1", nil)
	req.AddCookie(&http.Cookie{Name: "key", Value: "test-key"})

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("userID", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rec := httptest.NewRecorder()

	handler.GetAll(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response []entities.Note
	json.NewDecoder(rec.Body).Decode(&response)

	assert.Len(t, response, 2)
	assert.Equal(t, "Note 1", response[0].Title)
}

func TestNoteHandler_GetAll_InvalidUserID(t *testing.T) {
	mockService := new(MockNoteService)
	logger := zap.NewNop()
	handler := NewNoteHandler(mockService, logger)

	req := httptest.NewRequest(http.MethodGet, "/note/user/abc", nil)
	rec := httptest.NewRecorder()
	req.AddCookie(&http.Cookie{Name: "key", Value: "test-key"})
	req.Header.Set("Content-Type", "application/json")

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("userID", "abc")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.GetAll(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid user id")
}
