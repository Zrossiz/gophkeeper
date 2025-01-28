package service

import (
	"time"

	"github.com/Zrossiz/gophkeeper/internal/apperrors"
	"github.com/Zrossiz/gophkeeper/internal/config"
	"github.com/Zrossiz/gophkeeper/internal/dto"
	"github.com/Zrossiz/gophkeeper/internal/entities"
	"github.com/Zrossiz/gophkeeper/internal/utils"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	log    *zap.Logger
	dbUser UserStorage
	cfg    config.Config
}

type UserStorage interface {
	Create(body dto.UserDTO) error
	GetUserByUsername(username string) (*entities.User, error)
}

func NewUserService(
	dbUser UserStorage,
	cfg config.Config,
	logger *zap.Logger,
) *UserService {
	return &UserService{
		dbUser: dbUser,
		log:    logger,
		cfg:    cfg,
	}
}

func (u *UserService) Registration(registrationDTO dto.UserDTO) (*dto.GeneratedJwt, error) {
	hashedPassword, err := hashPassword(registrationDTO.Password, u.cfg.Cost)
	if err != nil {
		return nil, apperrors.ErrHashPassword
	}
	registrationDTO.Password = hashedPassword

	err = u.dbUser.Create(registrationDTO)
	if err != nil {
		u.log.Error(err.Error())
		return nil, apperrors.ErrDBQuery
	}

	createdUser, err := u.dbUser.GetUserByUsername(registrationDTO.Username)
	if err != nil {
		u.log.Error(err.Error())
		return nil, apperrors.ErrDBQuery
	}

	JWTAccessProps := utils.GenerateJWTProps{
		Secret:   []byte(u.cfg.AccessSecret),
		Exprires: time.Now().Add(15 * time.Minute),
		UserID:   int64(createdUser.ID),
		Username: createdUser.Username,
	}

	accessToken, err := utils.GenerateJWT(JWTAccessProps)
	if err != nil {
		u.log.Error(err.Error())
		return nil, apperrors.ErrJWTGeneration
	}

	JWTRefreshProps := utils.GenerateJWTProps{
		Secret:   []byte(u.cfg.RefreshSecret),
		Exprires: time.Now().Add(24 * 30 * time.Hour),
		UserID:   int64(createdUser.ID),
		Username: createdUser.Username,
	}

	refreshToken, err := utils.GenerateJWT(JWTRefreshProps)
	if err != nil {
		u.log.Error(err.Error())
		return nil, apperrors.ErrJWTGeneration
	}

	generatedTokens := dto.GeneratedJwt{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Hash:         createdUser.Password,
	}

	return &generatedTokens, nil
}

func (u *UserService) Login(loginDTO dto.UserDTO) (*dto.GeneratedJwt, error) {
	curUser, err := u.dbUser.GetUserByUsername(loginDTO.Username)
	if err != nil {
		u.log.Error(err.Error())
		return nil, apperrors.ErrDBQuery
	}

	if curUser == nil {
		return nil, apperrors.ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(curUser.Password), []byte(loginDTO.Password))
	if err != nil {
		return nil, apperrors.ErrInvalidPassword
	}

	JWTAccessProps := utils.GenerateJWTProps{
		Secret:   []byte(u.cfg.AccessSecret),
		Exprires: time.Now().Add(15 * time.Minute),
		UserID:   int64(curUser.ID),
		Username: curUser.Username,
	}

	accessToken, err := utils.GenerateJWT(JWTAccessProps)
	if err != nil {
		u.log.Error(err.Error())
		return nil, apperrors.ErrJWTGeneration
	}

	JWTRefreshProps := utils.GenerateJWTProps{
		Secret:   []byte(u.cfg.RefreshSecret),
		Exprires: time.Now().Add(24 * 30 * time.Hour),
		UserID:   int64(curUser.ID),
		Username: curUser.Username,
	}

	refreshToken, err := utils.GenerateJWT(JWTRefreshProps)
	if err != nil {
		u.log.Error(err.Error())
		return nil, apperrors.ErrJWTGeneration
	}

	generatedTokens := dto.GeneratedJwt{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Hash:         curUser.Password,
	}

	return &generatedTokens, nil
}

func hashPassword(password string, cost int) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
