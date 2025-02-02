package service

import (
	"github.com/Zrossiz/gophkeeper/internal/config"
	"go.uber.org/zap"
)

type Service struct {
	User     UserService
	LogoPass LogoPassService
	Binary   BinaryService
	Card     CardService
	Note     NoteService
}

type Storage struct {
	Binary   BinaryStorage
	User     UserStorage
	LogoPass LogoPassStorage
	Card     CardStorage
	Note     NoteStorage
}

type CryptoModule interface {
	Encrypt(plaintext, key string) (string, error)
	Decrypt(encryptedText, key string) (string, error)
	GenerateSecretPhrase(txt string) string
	EncryptBinaryData(plaintext []byte, key string) ([]byte, error)
	DecryptBinaryData(encryptedData []byte, key string) ([]byte, error)
}

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
