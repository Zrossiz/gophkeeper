package handler

import "go.uber.org/zap"

type Handler struct {
	User     UserHandler
	LogoPass LogoPassHandler
	Binary   BinaryHandler
	Card     CardHandler
}

type Service struct {
	User     UserService
	Card     CardService
	Binary   BinaryService
	LogoPass LogoPassService
}

func New(serv Service, logger *zap.Logger) *Handler {
	return &Handler{
		User:     *NewUserHandler(serv.User, logger),
		Binary:   *NewBinaryHandler(serv.Binary, logger),
		Card:     *NewCardHandler(serv.Card, logger),
		LogoPass: *NewLogoPassHandler(serv.LogoPass, logger),
	}
}
