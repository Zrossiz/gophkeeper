package service

import (
	"github.com/Zrossiz/gophkeeper/internal/dto"
	"github.com/Zrossiz/gophkeeper/internal/entities"
	"go.uber.org/zap"
)

type CardService struct {
	cardStorage  CardStorage
	cryptoModule CryptoModule
	log          *zap.Logger
}

type CardStorage interface {
	CreateCard(body dto.CreateCardDTO) error
	GetAllCardsByUserId(userID int64) ([]entities.Card, error)
	UpdateCard(cardID int64, body dto.UpdateCardDTO) error
}

func NewCardService(
	cardStorage CardStorage,
	cryptoModule CryptoModule,
	log *zap.Logger,
) *CardService {
	return &CardService{
		cardStorage:  cardStorage,
		cryptoModule: cryptoModule,
		log:          log,
	}
}

func (c *CardService) Create(body dto.CreateCardDTO) error {
	return c.cardStorage.CreateCard(body)
}

func (c *CardService) Update(cardID int64, body dto.UpdateCardDTO) error {
	return c.cardStorage.UpdateCard(cardID, body)
}

func (c *CardService) GetAll(userID int64) ([]entities.Card, error) {
	return c.cardStorage.GetAllCardsByUserId(userID)
}
