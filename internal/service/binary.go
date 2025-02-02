package service

import (
	"github.com/Zrossiz/gophkeeper/internal/dto"
	"github.com/Zrossiz/gophkeeper/internal/entities"
	"go.uber.org/zap"
)

type BinaryService struct {
	binaryStorage BinaryStorage
	cryptoModule  CryptoModule
	log           *zap.Logger
}

type BinaryStorage interface {
	Create(body dto.SetStorageBinaryDTO) error
	GetAllByUser(userID int64) ([]entities.BinaryData, error)
}

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

func (b *BinaryService) Create(body dto.CreateBinaryDTO) error {
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
	return b.binaryStorage.Create(binariesBody)
}

func (b *BinaryService) GetAll(userID int64, key string) ([]entities.BinaryData, error) {
	encryptedData, err := b.binaryStorage.GetAllByUser(userID)
	if err != nil {
		return nil, err
	}

	decryptedData := b.decryptBinaryArray(encryptedData, key)

	return decryptedData, nil
}

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
