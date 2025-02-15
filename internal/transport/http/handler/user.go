package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Zrossiz/gophkeeper/internal/apperrors"
	"github.com/Zrossiz/gophkeeper/internal/dto"
	"go.uber.org/zap"
)

type UserHandler struct {
	service UserService
	log     *zap.Logger
}

type UserService interface {
	Registration(ctx context.Context, registrationDTO dto.UserDTO) (*dto.GeneratedJwt, error)
	Login(ctx context.Context, loginDTO dto.UserDTO) (*dto.GeneratedJwt, error)
}

func NewUserHandler(serv UserService, log *zap.Logger) *UserHandler {
	return &UserHandler{service: serv}
}

// @Summary Регистрация пользователя
// @Description Создает нового пользователя в системе
// @Tags user
// @Accept  json
// @Produce  json
// @Param request body dto.UserDTO true "Данные пользователя"
// @Success 200
// @Failure 400
// @Failure 409
// @Failure 500
// @Router /api/user/register [post]
func (u *UserHandler) Registration(rw http.ResponseWriter, r *http.Request) {
	var registrationDTO dto.UserDTO

	err := json.NewDecoder(r.Body).Decode(&registrationDTO)
	if err != nil {
		http.Error(rw, "invalID request body", http.StatusBadRequest)
		return
	}

	if registrationDTO.Username == "" {
		http.Error(rw, "login can not be empty", http.StatusBadRequest)
		return
	}

	if registrationDTO.Password == "" {
		http.Error(rw, "password can not be empty", http.StatusBadRequest)
		return
	}

	generatedJwt, err := u.service.Registration(r.Context(), registrationDTO)
	if err != nil {
		switch err {
		case apperrors.ErrUserAlreadyExists:
			http.Error(rw, err.Error(), http.StatusConflict)
		case apperrors.ErrDBQuery:
			http.Error(rw, "internal server error", http.StatusInternalServerError)
		case apperrors.ErrHashPassword, apperrors.ErrJWTGeneration:
			http.Error(rw, "error processing request", http.StatusInternalServerError)
		default:
			u.log.Error(err.Error())
			http.Error(rw, "unknown error", http.StatusInternalServerError)
		}
		return
	}

	refreshTokenCokie := http.Cookie{
		Name:     "refreshtoken",
		Value:    generatedJwt.RefreshToken,
		Path:     "/",
		Expires:  time.Now().Add(2 * time.Hour * 24 * 30),
		HttpOnly: true,
		Secure:   false,
	}

	accessTokenCookie := http.Cookie{
		Name:     "accesstoken",
		Value:    generatedJwt.AccessToken,
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		Secure:   false,
	}

	keyCookie := http.Cookie{
		Name:     "key",
		Value:    generatedJwt.Hash,
		Path:     "/",
		Expires:  time.Now().Add(10000 * time.Hour),
		HttpOnly: true,
		Secure:   false,
	}

	http.SetCookie(rw, &refreshTokenCokie)
	http.SetCookie(rw, &accessTokenCookie)
	http.SetCookie(rw, &keyCookie)
	response := map[string]string{
		"hash": generatedJwt.Hash,
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(response); err != nil {
		http.Error(rw, "faile to encode response", http.StatusInternalServerError)
	}
}

// @Summary Авторизация пользователя
// @Description Аутентифицирует пользователя по логину и паролю
// @Tags user
// @Accept  json
// @Produce  json
// @Param request body dto.UserDTO true "Данные пользователя"
// @Success 200
// @Failure 400
// @Failure 401
// @Failure 500
// @Router /api/user/login [post]
func (u *UserHandler) Login(rw http.ResponseWriter, r *http.Request) {
	var loginDTO dto.UserDTO

	err := json.NewDecoder(r.Body).Decode(&loginDTO)
	if err != nil {
		http.Error(rw, "invalID request body", http.StatusBadRequest)
		return
	}

	if loginDTO.Username == "" {
		http.Error(rw, "login can not be empty", http.StatusBadRequest)
		return
	}

	if loginDTO.Password == "" {
		http.Error(rw, "password can not be empty", http.StatusBadRequest)
		return
	}

	generatedJwt, err := u.service.Login(r.Context(), loginDTO)
	if err != nil {
		switch err {
		case apperrors.ErrInvalidPassword:
			http.Error(rw, "unauthorized", http.StatusUnauthorized)
		case apperrors.ErrUserAlreadyExists:
			http.Error(rw, "user not found", http.StatusBadRequest)
		case apperrors.ErrDBQuery:
			http.Error(rw, "internal server error", http.StatusInternalServerError)
		case apperrors.ErrHashPassword, apperrors.ErrJWTGeneration:
			http.Error(rw, "error processing request", http.StatusInternalServerError)
		default:
			u.log.Error(err.Error())
			http.Error(rw, "unknown error", http.StatusInternalServerError)
		}
		return
	}

	refreshTokenCookie := http.Cookie{
		Name:     "refreshtoken",
		Value:    generatedJwt.RefreshToken,
		Path:     "/",
		Expires:  time.Now().Add(2 * time.Hour * 24 * 30),
		HttpOnly: true,
		Secure:   false,
	}

	accessTokenCookie := http.Cookie{
		Name:     "accesstoken",
		Value:    generatedJwt.AccessToken,
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		Secure:   false,
	}

	keyCookie := http.Cookie{
		Name:     "key",
		Value:    generatedJwt.Hash,
		Path:     "/",
		Expires:  time.Now().Add(10000 * time.Hour),
		HttpOnly: true,
		Secure:   false,
	}

	http.SetCookie(rw, &refreshTokenCookie)
	http.SetCookie(rw, &accessTokenCookie)
	http.SetCookie(rw, &keyCookie)

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
}
