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

func (l *LogoPassService) GetAll(userID int64, key string) ([]entities.LogoPassword, error) {
	items, err := l.logoPassDB.GetAllByUser(userID)
	if err != nil {
		return nil, err
	}

	decryptedData := l.decryptLogoPassArray(items, key)

	return decryptedData, nil
}

func (l *LogoPassService) decryptLogoPassArray(
	encryptedData []entities.LogoPassword,
	key string,
) []entities.LogoPassword {
	decryptedData := make([]entities.LogoPassword, 0, len(encryptedData))

	for i := 0; i < len(encryptedData); i++ {
		decryptedItem, err := l.decryptLogoPass(decryptedData[i], key)
		if err != nil {
			continue
		}

		decryptedData = append(decryptedData, *decryptedItem)
	}

	return decryptedData
}

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
