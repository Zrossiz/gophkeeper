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
	Create(body dto.CreateBinaryDTO) error
	Update(id int64, body dto.UpdateBinaryDTO) error
	GetAllByUser(userID int64) ([]entities.BinaryData, error)
}

func NewBinaryService(binaryStorage BinaryStorage, logger *zap.Logger) *BinaryService {
	return &BinaryService{
		binaryStorage: binaryStorage,
		log:           logger,
	}
}

func (b *BinaryService) Create(body dto.CreateBinaryDTO) error {
	return b.binaryStorage.Create(body)
}

func (b *BinaryService) Update(userID int64, body dto.UpdateBinaryDTO) error {
	return b.binaryStorage.Update(userID, body)
}

func (b *BinaryService) GetAll(userID int64) ([]entities.BinaryData, error) {
	return b.binaryStorage.GetAllByUser(userID)
}
