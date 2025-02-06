// Package service provides business logic for managing encrypted username-password storage.
package service

import (
	"fmt"

	"github.com/Zrossiz/gophkeeper/internal/dto"
	"github.com/Zrossiz/gophkeeper/internal/entities"
	"go.uber.org/zap"
)

// LogoPassService handles operations related to encrypted username-password storage.
type LogoPassService struct {
	logoPassDB   LogoPassStorage
	cryptoModule CryptoModule
	log          *zap.Logger
}

// LogoPassStorage defines an interface for storing, retrieving, and updating encrypted username-password data.
type LogoPassStorage interface {
	// CreateLogoPass stores an encrypted username-password entry.
	CreateLogoPass(body dto.CreateLogoPassDTO) error
	// GetAllByUser retrieves all encrypted username-password entries for a given user ID.
	GetAllByUser(userID int64) ([]entities.LogoPassword, error)
	// UpdateLogoPass updates an encrypted username-password entry for a given user ID.
	UpdateLogoPass(userID int64, body dto.UpdateLogoPassDTO) error
}

// NewLogoPassService creates a new instance of LogoPassService with the provided dependencies.
//
// Parameters:
//   - logoPassDB: An implementation of the LogoPassStorage interface for data persistence.
//   - cryptoModule: An implementation of CryptoModule for encryption and decryption.
//   - log: A structured logger (zap.Logger) for logging events.
//
// Returns:
//   - A pointer to a LogoPassService instance.
func NewLogoPassService(
	logoPassDB LogoPassStorage,
	cryptoModule CryptoModule,
	log *zap.Logger,
) *LogoPassService {
	return &LogoPassService{
		logoPassDB:   logoPassDB,
		cryptoModule: cryptoModule,
		log:          log,
	}
}

// Create encrypts and stores a username-password entry securely.
//
// Parameters:
//   - body: A dto.CreateLogoPassDTO containing username, password, and an encryption key.
//
// Returns:
//   - An error if encryption or storage fails.
func (l *LogoPassService) Create(body dto.CreateLogoPassDTO) error {
	encryptedUsername, err := l.cryptoModule.Encrypt(body.Username, body.Key)
	if err != nil {
		return err
	}

	encryptedPassword, err := l.cryptoModule.Encrypt(body.Password, body.Key)
	if err != nil {
		return err
	}

	body.Username = encryptedUsername
	body.Password = encryptedPassword

	err = l.logoPassDB.CreateLogoPass(body)
	if err != nil {
		return err
	}

	return nil
}

// Update encrypts and updates an existing username-password entry.
//
// Parameters:
//   - userID: The ID of the user whose data is being updated.
//   - body: A dto.UpdateLogoPassDTO containing updated username, password, and an encryption key.
//
// Returns:
//   - An error if encryption or update fails.
func (l *LogoPassService) Update(userID int64, body dto.UpdateLogoPassDTO) error {
	encryptedUsername, err := l.cryptoModule.Encrypt(body.Username, body.Key)
	if err != nil {
		return err
	}

	encryptedPassword, err := l.cryptoModule.Encrypt(body.Password, body.Key)
	if err != nil {
		return err
	}

	body.Username = encryptedUsername
	body.Password = encryptedPassword

	err = l.logoPassDB.UpdateLogoPass(userID, body)
	if err != nil {
		return err
	}

	return nil
}

// GetAll retrieves and decrypts all username-password entries for a given user.
//
// Parameters:
//   - userID: The ID of the user whose data is being retrieved.
//   - key: The encryption key required for decryption.
//
// Returns:
//   - A slice of decrypted entities.LogoPassword or an error if retrieval or decryption fails.
func (l *LogoPassService) GetAll(userID int64, key string) ([]entities.LogoPassword, error) {
	items, err := l.logoPassDB.GetAllByUser(userID)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("records not found")
	}

	decryptedData := l.decryptLogoPassArray(items, key)

	return decryptedData, nil
}

// decryptLogoPassArray decrypts an array of encrypted username-password entries.
//
// Parameters:
//   - encryptedData: A slice of encrypted entities.LogoPassword.
//   - key: The encryption key used for decryption.
//
// Returns:
//   - A slice of decrypted entities.LogoPassword.
func (l *LogoPassService) decryptLogoPassArray(
	encryptedData []entities.LogoPassword,
	key string,
) []entities.LogoPassword {
	decryptedData := make([]entities.LogoPassword, 0, len(encryptedData))

	for i := 0; i < len(encryptedData); i++ {
		decryptedItem, err := l.decryptLogoPass(encryptedData[i], key)
		if err != nil {
			continue
		}

		decryptedData = append(decryptedData, *decryptedItem)
	}

	return decryptedData
}

// decryptLogoPass decrypts a single encrypted username-password entry.
//
// Parameters:
//   - logopass: An encrypted entities.LogoPassword instance.
//   - key: The encryption key used for decryption.
//
// Returns:
//   - A pointer to a decrypted entities.LogoPassword or an error if decryption fails.
func (l *LogoPassService) decryptLogoPass(
	logopass entities.LogoPassword,
	key string,
) (*entities.LogoPassword, error) {
	decryptedLogin, err := l.cryptoModule.Decrypt(logopass.Username, key)
	if err != nil {
		return nil, err
	}

	decryptedPassword, err := l.cryptoModule.Decrypt(logopass.Password, key)
	if err != nil {
		return nil, err
	}

	logopass.Username = decryptedLogin
	logopass.Password = decryptedPassword

	return &logopass, nil
}
