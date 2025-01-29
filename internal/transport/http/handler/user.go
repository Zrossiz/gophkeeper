package handler

import (
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
	Registration(registrationDTO dto.UserDTO) (*dto.GeneratedJwt, error)
	Login(loginDTO dto.UserDTO) (*dto.GeneratedJwt, error)
}

func NewUserHandler(serv UserService, log *zap.Logger) *UserHandler {
	return &UserHandler{service: serv}
}

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

	generatedJwt, err := u.service.Registration(registrationDTO)
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

	http.SetCookie(rw, &refreshTokenCokie)
	http.SetCookie(rw, &accessTokenCookie)
	response := map[string]string{
		"message": "registration successful",
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(response)
}

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

	generatedJwt, err := u.service.Login(loginDTO)
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

	http.SetCookie(rw, &refreshTokenCookie)
	http.SetCookie(rw, &accessTokenCookie)
	response := map[string]string{
		"message": "login successful",
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	err = json.NewEncoder(rw).Encode(response)
	if err != nil {
		http.Error(rw, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
