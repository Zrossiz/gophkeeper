package handler

import "go.uber.org/zap"

type Handler struct {
	User     UserHandler
	LogoPass LogoPassHandler
	Binary   BinaryHandler
	Card     CardHandler
	Note     NoteHandler
}

type Service struct {
	User     UserService
	Card     CardService
	Binary   BinaryService
	LogoPass LogoPassService
	Note     NoteService
}

func New(serv Service, logger *zap.Logger) *Handler {
	return &Handler{
		User:     *NewUserHandler(serv.User, logger),
		Binary:   *NewBinaryHandler(serv.Binary, logger),
		Card:     *NewCardHandler(serv.Card, logger),
		LogoPass: *NewLogoPassHandler(serv.LogoPass, logger),
		Note:     *NewNoteHandler(serv.Note, logger),
	}
}
