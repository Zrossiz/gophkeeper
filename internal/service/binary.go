package service

import (
	"github.com/Zrossiz/gophkeeper/internal/dto"
	"github.com/Zrossiz/gophkeeper/internal/entities"
	"go.uber.org/zap"
)

type BinaryService struct {
	binaryStorage BinaryStorage
	log           *zap.Logger
}

type BinaryStorage interface {
	Create(body dto.SetStorageBinaryDTO) error
	Update(body dto.SetStorageBinaryDTO) error
	GetAllByUser(userID int64) ([]entities.BinaryData, error)
}

func NewBinaryService(binaryStorage BinaryStorage, logger *zap.Logger) *BinaryService {
	return &BinaryService{
		binaryStorage: binaryStorage,
		log:           logger,
	}
}

func (b *BinaryService) Create(body dto.CreateBinaryDTO) error {
	binariesBody := dto.SetStorageBinaryDTO{
		UserID: body.UserId,
		Title:  body.Title,
		Data:   []byte(body.Data),
	}
	return b.binaryStorage.Create(binariesBody)
}

func (b *BinaryService) Update(userID int, body dto.UpdateBinaryDTO) error {
	binariesBody := dto.SetStorageBinaryDTO{
		UserID: userID,
		Title:  body.Title,
		Data:   []byte(body.Data),
	}
	return b.binaryStorage.Update(binariesBody)
}

func (b *BinaryService) GetAll(userID int64) ([]entities.BinaryData, error) {
	return b.binaryStorage.GetAllByUser(userID)
}
