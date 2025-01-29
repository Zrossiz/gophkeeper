package service

import (
	"github.com/Zrossiz/gophkeeper/internal/dto"
	"github.com/Zrossiz/gophkeeper/internal/entities"
	"go.uber.org/zap"
)

type LogoPassService struct {
	logoPassDB   LogoPassStorage
	cryptoModule CryptoModule
	log          *zap.Logger
}

type LogoPassStorage interface {
	CreateLogoPass(body dto.CreateLogoPassDTO) error
	GetAllByUser(userID int64) ([]entities.LogoPassword, error)
	UpdateLogoPass(userID int64, body dto.UpdateLogoPassDTO) error
}

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

func (l *LogoPassService) Create(body dto.CreateLogoPassDTO) error {
	err := l.logoPassDB.CreateLogoPass(body)
	if err != nil {
		return err
	}

	return nil
}

func (l *LogoPassService) Update(userID int64, body dto.UpdateLogoPassDTO) error {
	err := l.logoPassDB.UpdateLogoPass(userID, body)
	if err != nil {
		return err
	}

	return nil
}

func (l *LogoPassService) GetAll(userID int64) ([]entities.LogoPassword, error) {
	items, err := l.logoPassDB.GetAllByUser(userID)
	if err != nil {
		return nil, err
	}

	return items, nil
}
