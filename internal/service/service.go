package service

type Service struct {
	User     UserService
	LogoPass LogoPassService
	Binary   BinaryService
	Card     CardService
}

func New() *Service {
	return &Service{
		User:     *NewUserService(),
		Binary:   *NewBinaryService(),
		Card:     *NewCardService(),
		LogoPass: *NewLogoPassService(),
	}
}
