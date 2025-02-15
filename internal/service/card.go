// Package service provides business logic for managing encrypted card data storage.
package service

import (
	"context"
	"fmt"

	"github.com/Zrossiz/gophkeeper/internal/dto"
	"github.com/Zrossiz/gophkeeper/internal/entities"
	"go.uber.org/zap"
)

// CardService manages operations related to encrypted card data storage, including encryption and decryption.
type CardService struct {
	cardStorage  CardStorage
	cryptoModule CryptoModule
	log          *zap.Logger
}

// CardStorage defines an interface for storing, retrieving, and updating encrypted card data.
type CardStorage interface {
	// CreateCard stores an encrypted card in the database.
	CreateCard(ctx context.Context, body dto.CreateCardDTO) error
	// GetAllCardsByUserId retrieves all encrypted cards associated with a given user ID.
	GetAllCardsByUserId(ctx context.Context, userID int64) ([]entities.Card, error)
	// UpdateCard updates the encrypted card details for a specific card ID.
	UpdateCard(ctx context.Context, cardID int64, body dto.UpdateCardDTO) error
}

// NewCardService creates a new instance of CardService with the provided dependencies.
//
// Parameters:
//   - cardStorage: An implementation of the CardStorage interface for data persistence.
//   - cryptoModule: An implementation of CryptoModule for encryption and decryption.
//   - log: A structured logger (zap.Logger) for logging events.
//
// Returns:
//   - A pointer to a CardService instance.
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

// Create encrypts card data and stores it securely.
//
// Parameters:
//   - body: A dto.CreateCardDTO containing card details and an encryption key.
//
// Returns:
//   - An error if encryption or storage fails.
func (c *CardService) Create(ctx context.Context, body dto.CreateCardDTO) error {
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

	encryptedCardHolderName, err := c.cryptoModule.Encrypt(body.CardHolderName, body.Key)
	if err != nil {
		return err
	}

	body.Num = encryptedNum
	body.CVV = encryptedCVV
	body.ExpDate = encryptedExpDate
	body.CardHolderName = encryptedCardHolderName

	return c.cardStorage.CreateCard(ctx, body)
}

// Update encrypts updated card data and stores it securely.
//
// Parameters:
//   - cardID: The ID of the card to be updated.
//   - body: A dto.UpdateCardDTO containing updated card details and an encryption key.
//
// Returns:
//   - An error if encryption or storage fails.
func (c *CardService) Update(ctx context.Context, cardID int64, body dto.UpdateCardDTO) error {
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

	encryptedCardHolderName, err := c.cryptoModule.Encrypt(body.CardHolderName, body.Key)
	if err != nil {
		return err
	}

	body.Num = encryptedNum
	body.CVV = encryptedCVV
	body.ExpDate = encryptedExpDate
	body.CardHolderName = encryptedCardHolderName

	return c.cardStorage.UpdateCard(ctx, cardID, body)
}

// GetAll retrieves and decrypts all card data for a given user.
//
// Parameters:
//   - userID: The ID of the user whose card data is being retrieved.
//   - key: The encryption key required for decryption.
//
// Returns:
//   - A slice of decrypted entities.Card or an error if retrieval or decryption fails.
func (c *CardService) GetAll(ctx context.Context, userID int64, key string) ([]entities.Card, error) {
	encryptedData, err := c.cardStorage.GetAllCardsByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(encryptedData) == 0 {
		return nil, fmt.Errorf("records not found")
	}

	decryptedData := c.decryptCardArray(encryptedData, key)
	return decryptedData, nil
}

// decryptCardArray decrypts an array of encrypted card data.
//
// Parameters:
//   - cards: A slice of encrypted entities.Card.
//   - key: The encryption key used for decryption.
//
// Returns:
//   - A slice of decrypted entities.Card.
func (c *CardService) decryptCardArray(cards []entities.Card, key string) []entities.Card {
	decryptedData := make([]entities.Card, 0, len(cards))

	for i := 0; i < len(cards); i++ {
		decryptedCard, err := c.decryptCard(cards[i], key)
		if err != nil {
			continue
		}
		decryptedData = append(decryptedData, *decryptedCard)
	}

	return decryptedData
}

// decryptCard decrypts a single encrypted card entry.
//
// Parameters:
//   - card: An encrypted entities.Card instance.
//   - key: The encryption key used for decryption.
//
// Returns:
//   - A pointer to a decrypted entities.Card or an error if decryption fails.
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

	decryptedCardHolderName, err := c.cryptoModule.Decrypt(card.CardHolderName, key)
	if err != nil {
		return nil, err
	}

	card.Number = decryptedNum
	card.ExpDate = decryptedExpDate
	card.CVV = decryptedCVV
	card.CardHolderName = decryptedCardHolderName

	return &card, nil
}
