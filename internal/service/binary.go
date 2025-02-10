// Package service provides business logic for managing binary data storage and encryption.
package service

import (
	"context"

	"github.com/Zrossiz/gophkeeper/internal/dto"
	"github.com/Zrossiz/gophkeeper/internal/entities"
	"go.uber.org/zap"
)

// BinaryService handles operations related to binary data storage, including encryption and decryption.
type BinaryService struct {
	binaryStorage BinaryStorage
	cryptoModule  CryptoModule
	log           *zap.Logger
}

// BinaryStorage defines an interface for storing and retrieving encrypted binary data.
type BinaryStorage interface {
	// Create stores encrypted binary data.
	Create(ctx context.Context, body dto.SetStorageBinaryDTO) error
	// GetAllByUser retrieves all binary data associated with a given user.
	GetAllByUser(ctx context.Context, userID int64) ([]entities.BinaryData, error)
}

// NewBinaryService creates a new instance of BinaryService with the provided dependencies.
//
// Parameters:
//   - binaryStorage: An implementation of the BinaryStorage interface for data persistence.
//   - cryptoModule: An implementation of CryptoModule for encryption and decryption.
//   - logger: A structured logger (zap.Logger) for logging events.
//
// Returns:
//   - A pointer to a BinaryService instance.
func NewBinaryService(
	binaryStorage BinaryStorage,
	cryptoModule CryptoModule,
	logger *zap.Logger,
) *BinaryService {
	return &BinaryService{
		binaryStorage: binaryStorage,
		cryptoModule:  cryptoModule,
		log:           logger,
	}
}

// Create encrypts binary data and stores it securely.
//
// Parameters:
//   - body: A dto.CreateBinaryDTO containing user ID, title, data, and encryption key.
//
// Returns:
//   - An error if encryption or storage fails.
func (b *BinaryService) Create(ctx context.Context, body dto.CreateBinaryDTO) error {
	encryptedTitle, err := b.cryptoModule.Encrypt(body.Title, body.Key)
	if err != nil {
		return err
	}

	encryptedBody, err := b.cryptoModule.EncryptBinaryData(body.Data, body.Key)
	if err != nil {
		return err
	}

	binariesBody := dto.SetStorageBinaryDTO{
		UserID: body.UserID,
		Title:  encryptedTitle,
		Data:   encryptedBody,
	}
	return b.binaryStorage.Create(ctx, binariesBody)
}

// GetAll retrieves and decrypts all binary data for a given user.
//
// Parameters:
//   - userID: The ID of the user whose data is being retrieved.
//   - key: The encryption key required for decryption.
//
// Returns:
//   - A slice of decrypted entities.BinaryData or an error if retrieval or decryption fails.
func (b *BinaryService) GetAll(ctx context.Context, userID int64, key string) ([]entities.BinaryData, error) {
	encryptedData, err := b.binaryStorage.GetAllByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	decryptedData := b.decryptBinaryArray(encryptedData, key)

	return decryptedData, nil
}

// decryptBinaryArray decrypts an array of encrypted binary data.
//
// Parameters:
//   - encryptedData: A slice of encrypted entities.BinaryData.
//   - key: The encryption key used for decryption.
//
// Returns:
//   - A slice of decrypted entities.BinaryData.
func (b *BinaryService) decryptBinaryArray(
	encryptedData []entities.BinaryData,
	key string,
) []entities.BinaryData {
	decryptedData := make([]entities.BinaryData, 0, len(encryptedData))

	for i := 0; i < len(encryptedData); i++ {
		decryptedItem, err := b.decryptBinary(encryptedData[i], key)
		if err != nil {
			continue
		}

		decryptedData = append(decryptedData, *decryptedItem)
	}

	return decryptedData
}

// decryptBinary decrypts a single encrypted binary data entry.
//
// Parameters:
//   - encryptedData: An encrypted entities.BinaryData instance.
//   - key: The encryption key used for decryption.
//
// Returns:
//   - A pointer to a decrypted entities.BinaryData or an error if decryption fails.
func (b *BinaryService) decryptBinary(
	encryptedData entities.BinaryData,
	key string,
) (*entities.BinaryData, error) {
	decryptedTitle, err := b.cryptoModule.Decrypt(encryptedData.Title, key)
	if err != nil {
		return nil, err
	}

	encryptedData.Title = decryptedTitle

	return &encryptedData, nil
}
