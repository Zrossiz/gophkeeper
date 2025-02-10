package service

import (
	"context"
	"testing"

	"github.com/Zrossiz/gophkeeper/internal/dto"
	"github.com/Zrossiz/gophkeeper/internal/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockNoteStorage struct {
	mock.Mock
}

func (m *MockNoteStorage) Create(ctx context.Context, body dto.CreateNoteDTO) error {
	args := m.Called(body)
	return args.Error(0)
}

func (m *MockNoteStorage) Update(ctx context.Context, noteID int, body dto.UpdateNoteDTO) error {
	args := m.Called(noteID, body)
	return args.Error(0)
}

func (m *MockNoteStorage) GetAllByUser(ctx context.Context, userID int) ([]entities.Note, error) {
	args := m.Called(userID)
	return args.Get(0).([]entities.Note), args.Error(1)
}

func TestCreateNote(t *testing.T) {
	mockStorage := new(MockNoteStorage)
	mockCrypto := new(MockCryptoModule)
	logger := zap.NewNop()

	service := NewNoteService(mockStorage, mockCrypto, logger)

	noteDTO := dto.CreateNoteDTO{
		Title:    "My Note",
		TextData: "This is a test note",
		Key:      "secret",
	}

	mockCrypto.On("Encrypt", "My Note", "secret").Return("enc_title", nil)
	mockCrypto.On("Encrypt", "This is a test note", "secret").Return("enc_text", nil)

	mockStorage.On("Create", mock.Anything).Return(nil)

	err := service.Create(context.Background(), noteDTO)

	assert.NoError(t, err)
	mockCrypto.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
}

func TestUpdateNote(t *testing.T) {
	mockStorage := new(MockNoteStorage)
	mockCrypto := new(MockCryptoModule)
	logger := zap.NewNop()

	service := NewNoteService(mockStorage, mockCrypto, logger)

	noteDTO := dto.UpdateNoteDTO{
		Title:    "Updated Title",
		TextData: "Updated text",
		Key:      "secret",
	}

	noteID := 1

	mockCrypto.On("Encrypt", "Updated Title", "secret").Return("enc_title", nil)
	mockCrypto.On("Encrypt", "Updated text", "secret").Return("enc_text", nil)

	mockStorage.On("Update", noteID, mock.Anything).Return(nil)

	err := service.Update(context.Background(), noteID, noteDTO)

	assert.NoError(t, err)
	mockCrypto.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
}

func TestGetAllNotes_Success(t *testing.T) {
	mockStorage := new(MockNoteStorage)
	mockCrypto := new(MockCryptoModule)
	logger := zap.NewNop()

	service := NewNoteService(mockStorage, mockCrypto, logger)

	userID := 1
	encryptionKey := "secret"

	encryptedNotes := []entities.Note{
		{Title: "enc_title1", TextData: "enc_text1"},
		{Title: "enc_title2", TextData: "enc_text2"},
	}

	mockStorage.On("GetAllByUser", userID).Return(encryptedNotes, nil)

	mockCrypto.On("Decrypt", "enc_title1", encryptionKey).Return("Title 1", nil)
	mockCrypto.On("Decrypt", "enc_text1", encryptionKey).Return("Text 1", nil)
	mockCrypto.On("Decrypt", "enc_title2", encryptionKey).Return("Title 2", nil)
	mockCrypto.On("Decrypt", "enc_text2", encryptionKey).Return("Text 2", nil)

	notes, err := service.GetAll(context.Background(), userID, encryptionKey)

	assert.NoError(t, err)
	assert.Len(t, notes, 2)
	assert.Equal(t, "Title 1", notes[0].Title)
	assert.Equal(t, "Text 1", notes[0].TextData)
	assert.Equal(t, "Title 2", notes[1].Title)
	assert.Equal(t, "Text 2", notes[1].TextData)

	mockStorage.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
}

func TestGetAllNotes_NotFound(t *testing.T) {
	mockStorage := new(MockNoteStorage)
	mockCrypto := new(MockCryptoModule)
	logger := zap.NewNop()

	service := NewNoteService(mockStorage, mockCrypto, logger)

	mockStorage.On("GetAllByUser", 1).Return([]entities.Note{}, nil)

	notes, err := service.GetAll(context.Background(), 1, "secret")

	assert.NoError(t, err)
	assert.Empty(t, notes)

	mockStorage.AssertExpectations(t)
}
