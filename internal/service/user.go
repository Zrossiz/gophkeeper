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

// UserService handles user authentication, registration, and JWT generation.
type UserService struct {
	log          *zap.Logger   // Logger for structured logging.
	dbUser       UserStorage   // Database storage interface for user data.
	cfg          config.Config // Application configuration.
	cryptoModule CryptoModule  // Cryptographic module for password security.
}

// UserStorage defines database operations related to user management.
type UserStorage interface {
	// Create saves a new user record in the database.
	Create(body dto.UserDTO) error
	// GetUserByUsername retrieves a user by their username.
	GetUserByUsername(username string) (*entities.User, error)
}

// NewUserService initializes and returns a new UserService instance.
//
// Parameters:
//   - dbUser: Implementation of UserStorage interface.
//   - cryptoModule: Cryptographic module for password security.
//   - cfg: Application configuration.
//   - logger: Structured logger (zap.Logger).
//
// Returns:
//   - A pointer to a fully initialized UserService instance.
func NewUserService(
	dbUser UserStorage,
	cryptoModule CryptoModule,
	cfg config.Config,
	logger *zap.Logger,
) *UserService {
	return &UserService{
		dbUser:       dbUser,
		cryptoModule: cryptoModule,
		log:          logger,
		cfg:          cfg,
	}
}

// Registration registers a new user, hashes their password, and generates JWT tokens.
//
// Parameters:
//   - registrationDTO: Contains user registration details (username, password).
//
// Returns:
//   - A pointer to GeneratedJwt struct containing access and refresh tokens.
//   - An error if user creation fails or token generation fails.
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
		Hash:         u.cryptoModule.GenerateSecretPhrase(createdUser.Password),
	}

	return &generatedTokens, nil
}

// Login authenticates a user and generates JWT tokens.
//
// Parameters:
//   - loginDTO: Contains user login details (username, password).
//
// Returns:
//   - A pointer to GeneratedJwt struct containing access and refresh tokens.
//   - An error if authentication fails.
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
	}

	return &generatedTokens, nil
}

// hashPassword hashes a given password using bcrypt.
//
// Parameters:
//   - password: The plain text password to be hashed.
//   - cost: The bcrypt cost factor (determines computation complexity).
//
// Returns:
//   - A hashed password as a string.
//   - An error if hashing fails.
func hashPassword(password string, cost int) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
