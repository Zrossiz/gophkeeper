package service

import (
	"testing"

	"github.com/Zrossiz/gophkeeper/internal/dto"
	"github.com/Zrossiz/gophkeeper/internal/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockCardStorage struct {
	mock.Mock
}

func (m *MockCardStorage) CreateCard(body dto.CreateCardDTO) error {
	args := m.Called(body)
	return args.Error(0)
}

func (m *MockCardStorage) GetAllCardsByUserId(userID int64) ([]entities.Card, error) {
	args := m.Called(userID)
	return args.Get(0).([]entities.Card), args.Error(1)
}

func (m *MockCardStorage) UpdateCard(cardID int64, body dto.UpdateCardDTO) error {
	args := m.Called(cardID, body)
	return args.Error(0)
}

type MockCryptoModule struct {
	mock.Mock
}

func (m *MockCryptoModule) Encrypt(data, key string) (string, error) {
	args := m.Called(data, key)
	return args.String(0), args.Error(1)
}

func (m *MockCryptoModule) Decrypt(data, key string) (string, error) {
	args := m.Called(data, key)
	return args.String(0), args.Error(1)
}

func (m *MockCryptoModule) EncryptBinaryData(data []byte, key string) ([]byte, error) {
	args := m.Called(data, key)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockCryptoModule) DecryptBinaryData(data []byte, key string) ([]byte, error) {
	args := m.Called(data, key)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockCryptoModule) GenerateSecretPhrase(txt string) string {
	args := m.Called(txt)
	return args.String(0)
}

func TestCreateCard(t *testing.T) {
	mockStorage := new(MockCardStorage)
	mockCrypto := new(MockCryptoModule)
	logger := zap.NewNop()

	service := NewCardService(mockStorage, mockCrypto, logger)

	cardDTO := dto.CreateCardDTO{
		Num:            "1234 5678 9101 1121",
		CVV:            "123",
		ExpDate:        "12/25",
		CardHolderName: "John Doe",
		Key:            "secret",
	}

	mockCrypto.On("Encrypt", "1234 5678 9101 1121", "secret").Return("encrypted_num", nil)
	mockCrypto.On("Encrypt", "123", "secret").Return("encrypted_cvv", nil)
	mockCrypto.On("Encrypt", "12/25", "secret").Return("encrypted_exp", nil)
	mockCrypto.On("Encrypt", "John Doe", "secret").Return("encrypted_name", nil)

	mockStorage.On("CreateCard", mock.Anything).Return(nil)

	err := service.Create(cardDTO)

	assert.NoError(t, err)
	mockCrypto.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
}

func TestUpdateCard(t *testing.T) {
	mockStorage := new(MockCardStorage)
	mockCrypto := new(MockCryptoModule)
	logger := zap.NewNop()

	service := NewCardService(mockStorage, mockCrypto, logger)

	cardDTO := dto.UpdateCardDTO{
		Num:            "9876 5432 1098 7654",
		CVV:            "321",
		ExpDate:        "11/24",
		CardHolderName: "Jane Doe",
		Key:            "secret",
	}

	mockCrypto.On("Encrypt", "9876 5432 1098 7654", "secret").Return("encrypted_num", nil)
	mockCrypto.On("Encrypt", "321", "secret").Return("encrypted_cvv", nil)
	mockCrypto.On("Encrypt", "11/24", "secret").Return("encrypted_exp", nil)
	mockCrypto.On("Encrypt", "Jane Doe", "secret").Return("encrypted_name", nil)

	mockStorage.On("UpdateCard", int64(1), mock.Anything).Return(nil)

	err := service.Update(1, cardDTO)

	assert.NoError(t, err)
	mockCrypto.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
}

func TestGetAllCards_Success(t *testing.T) {
	mockStorage := new(MockCardStorage)
	mockCrypto := new(MockCryptoModule)
	logger := zap.NewNop()

	service := NewCardService(mockStorage, mockCrypto, logger)

	encryptedCards := []entities.Card{
		{Number: "enc_1", CVV: "enc_2", ExpDate: "enc_3", CardHolderName: "enc_4"},
	}

	mockStorage.On("GetAllCardsByUserId", int64(1)).Return(encryptedCards, nil)

	mockCrypto.On("Decrypt", "enc_1", "secret").Return("1234 5678 9101 1121", nil)
	mockCrypto.On("Decrypt", "enc_2", "secret").Return("123", nil)
	mockCrypto.On("Decrypt", "enc_3", "secret").Return("12/25", nil)
	mockCrypto.On("Decrypt", "enc_4", "secret").Return("John Doe", nil)

	cards, err := service.GetAll(1, "secret")

	assert.NoError(t, err)
	assert.Len(t, cards, 1)
	assert.Equal(t, "1234 5678 9101 1121", cards[0].Number)
	assert.Equal(t, "123", cards[0].CVV)
	assert.Equal(t, "12/25", cards[0].ExpDate)
	assert.Equal(t, "John Doe", cards[0].CardHolderName)

	mockStorage.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
}

func TestGetAllCards_NotFound(t *testing.T) {
	mockStorage := new(MockCardStorage)
	mockCrypto := new(MockCryptoModule)
	logger := zap.NewNop()

	service := NewCardService(mockStorage, mockCrypto, logger)

	mockStorage.On("GetAllCardsByUserId", int64(1)).Return([]entities.Card{}, nil)

	cards, err := service.GetAll(1, "secret")

	assert.Error(t, err)
	assert.Nil(t, cards)
	assert.Equal(t, "records not found", err.Error())

	mockStorage.AssertExpectations(t)
}
