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
}

type Storage struct {
	Binary   BinaryStorage
	User     UserStorage
	LogoPass LogoPassStorage
	Card     CardStorage
}

func New(
	store Storage,
	cfg config.Config,
	logger *zap.Logger,
) *Service {
	return &Service{
		User:     *NewUserService(store.User, cfg, logger),
		Binary:   *NewBinaryService(store.Binary, logger),
		Card:     *NewCardService(store.Card, logger),
		LogoPass: *NewLogoPassService(store.LogoPass, logger),
	}
}
