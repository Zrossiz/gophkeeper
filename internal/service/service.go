// Package service provides a collection of services for managing encrypted user data.
package service

import (
	"github.com/Zrossiz/gophkeeper/internal/config"
	"go.uber.org/zap"
)

// Service aggregates all individual services responsible for managing different types of data.
type Service struct {
	User     UserService     // Handles user authentication and management.
	LogoPass LogoPassService // Manages encrypted login-password storage.
	Binary   BinaryService   // Manages encrypted binary file storage.
	Card     CardService     // Handles encrypted card data storage.
	Note     NoteService     // Manages encrypted note storage.
}

// Storage defines interfaces for data persistence layers corresponding to different services.
type Storage struct {
	Binary   BinaryStorage   // Interface for binary data storage operations.
	User     UserStorage     // Interface for user data storage operations.
	LogoPass LogoPassStorage // Interface for login-password storage operations.
	Card     CardStorage     // Interface for card data storage operations.
	Note     NoteStorage     // Interface for note storage operations.
}

// CryptoModule defines an interface for cryptographic operations used throughout the services.
type CryptoModule interface {
	// Encrypt encrypts the given plaintext using the provided key.
	Encrypt(plaintext, key string) (string, error)
	// Decrypt decrypts the given encrypted text using the provided key.
	Decrypt(encryptedText, key string) (string, error)
	// GenerateSecretPhrase generates a secret phrase from the provided text.
	GenerateSecretPhrase(txt string) string
	// EncryptBinaryData encrypts binary data using the provided key.
	EncryptBinaryData(plaintext []byte, key string) ([]byte, error)
	// DecryptBinaryData decrypts binary data using the provided key.
	DecryptBinaryData(encryptedData []byte, key string) ([]byte, error)
}

// New initializes and returns a Service instance with all dependencies injected.
//
// Parameters:
//   - store: A Storage instance containing implementations of various storage interfaces.
//   - cfg: A configuration object containing application settings.
//   - cryptoModule: An implementation of the CryptoModule interface for encryption and decryption.
//   - logger: A structured logger (zap.Logger) for logging events.
//
// Returns:
//   - A pointer to a fully initialized Service instance.
func New(
	store Storage,
	cfg config.Config,
	cryptoModule CryptoModule,
	logger *zap.Logger,
) *Service {
	return &Service{
		User:     *NewUserService(store.User, cryptoModule, cfg, logger),
		Binary:   *NewBinaryService(store.Binary, cryptoModule, logger),
		Card:     *NewCardService(store.Card, cryptoModule, logger),
		LogoPass: *NewLogoPassService(store.LogoPass, cryptoModule, logger),
		Note:     *NewNoteService(store.Note, cryptoModule, logger),
	}
}
