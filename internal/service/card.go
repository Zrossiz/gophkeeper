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
	encryptedNum, err := c.cryptoModule.Encrypt(body.Num, body.Key)
	if err != nil {
		return err
	}

	encryptedCVV, err := c.cryptoModule.Encrypt(body.CVV, body.Key)
	if err != nil {
		return err
	}

	encryptedExpDate, err := c.cryptoModule.Encrypt(body.ExpDate, body.Key)
	if err != nil {
		return err
	}

	body.Num = encryptedNum
	body.CVV = encryptedCVV
	body.ExpDate = encryptedExpDate

	return c.cardStorage.CreateCard(body)
}

func (c *CardService) Update(cardID int64, body dto.UpdateCardDTO) error {
	encryptedNum, err := c.cryptoModule.Encrypt(body.Num, body.Key)
	if err != nil {
		return err
	}

	encryptedCVV, err := c.cryptoModule.Encrypt(body.CVV, body.Key)
	if err != nil {
		return err
	}

	encryptedExpDate, err := c.cryptoModule.Encrypt(body.ExpDate, body.Key)
	if err != nil {
		return err
	}

	body.Num = encryptedNum
	body.CVV = encryptedCVV
	body.ExpDate = encryptedExpDate

	return c.cardStorage.UpdateCard(cardID, body)
}

func (c *CardService) GetAll(userID int64, key string) ([]entities.Card, error) {
	encryptedData, err := c.cardStorage.GetAllCardsByUserId(userID)
	if err != nil {
		return nil, err
	}

	decryptedData := c.decryptCardArray(encryptedData, key)

	return decryptedData, nil
}

func (c *CardService) decryptCardArray(card []entities.Card, key string) []entities.Card {
	decryptedData := make([]entities.Card, 0, len(card))

	for i := 0; i < len(card); i++ {
		decryptedCard, err := c.decryptCard(card[i], key)
		if err != nil {
			continue
		}
		decryptedData = append(decryptedData, *decryptedCard)
	}

	return decryptedData
}

func (c *CardService) decryptCard(card entities.Card, key string) (*entities.Card, error) {
	decryptedNum, err := c.cryptoModule.Decrypt(card.Number, key)
	if err != nil {
		return nil, err
	}

	decryptedCVV, err := c.cryptoModule.Decrypt(card.CVV, key)
	if err != nil {
		return nil, err
	}

	decryptedExpDate, err := c.cryptoModule.Decrypt(card.ExpDate, key)
	if err != nil {
		return nil, err
	}

	card.Number = decryptedNum
	card.ExpDate = decryptedExpDate
	card.CVV = decryptedCVV

	return &card, nil
}
