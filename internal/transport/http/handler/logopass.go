package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Zrossiz/gophkeeper/internal/apperrors"
	"github.com/Zrossiz/gophkeeper/internal/dto"
	"github.com/Zrossiz/gophkeeper/internal/entities"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type LogoPassHandler struct {
	service LogoPassService
	log     *zap.Logger
}

type LogoPassService interface {
	Create(body dto.CreateLogoPassDTO) error
	Update(userID int64, body dto.UpdateLogoPassDTO) error
	GetAll(userID int64, key string) ([]entities.LogoPassword, error)
}

func NewLogoPassHandler(service LogoPassService, logger *zap.Logger) *LogoPassHandler {
	return &LogoPassHandler{
		service: service,
		log:     logger,
	}
}

// @Summary Создать логин-пароль
// @Description Создает новую запись логина и пароля пользователя
// @Tags logopass
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer токен" default(Bearer {token})
// @Param body body dto.CreateLogoPassDTO true "Данные для создания логина и пароля"
// @Success 201 {string} string "Created"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /logo-pass [post]
// @Security BearerAuth
func (l *LogoPassHandler) Create(rw http.ResponseWriter, r *http.Request) {
	key, err := r.Cookie("key")
	if err != nil {
		http.Error(rw, "key not found", http.StatusBadRequest)
		return
	}

	var body dto.CreateLogoPassDTO
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(rw, apperrors.ErrInvalidRequestBody, http.StatusBadRequest)
		return
	}

	body.Key = key.Value

	err = l.service.Create(body)
	if err != nil {
		l.log.Sugar().Errorf("create logo pass error: %v", err)
		http.Error(rw, apperrors.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

// @Summary Обновить логин-пароль
// @Description Обновляет существующую запись логина и пароля
// @Tags logopass
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer токен" default(Bearer {token})
// @Param logoPassID path int true "ID логина-пароля"
// @Param body body dto.UpdateLogoPassDTO true "Данные для обновления"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /logo-pass/{logoPassID} [put]
// @Security BearerAuth
func (l *LogoPassHandler) Update(rw http.ResponseWriter, r *http.Request) {
	key, err := r.Cookie("key")
	if err != nil {
		http.Error(rw, "key not found", http.StatusBadRequest)
		return
	}

	var body dto.UpdateLogoPassDTO
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(rw, apperrors.ErrInvalidRequestBody, http.StatusBadRequest)
		return
	}

	logoPassID := chi.URLParam(r, "logoPassID")
	intLogoPassID, err := strconv.Atoi(logoPassID)
	if err != nil {
		http.Error(rw, "invalid logo pass id ", http.StatusBadRequest)
		return
	}

	body.Key = key.Value

	err = l.service.Update(int64(intLogoPassID), body)
	if err != nil {
		l.log.Sugar().Errorf("update logo pass error: %v", err)
		http.Error(rw, apperrors.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

// @Summary Получить все логин-пароли пользователя
// @Description Возвращает список всех сохраненных логинов и паролей
// @Tags logopass
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer токен" default(Bearer {token})
// @Param userID path int true "ID пользователя"
// @Success 200 {array} entities.LogoPassword "Список логинов и паролей"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /logo-pass/user/{userID} [get]
// @Security BearerAuth
func (l *LogoPassHandler) GetAll(rw http.ResponseWriter, r *http.Request) {
	key, err := r.Cookie("key")
	if err != nil {
		http.Error(rw, "key not found", http.StatusBadRequest)
		return
	}

	userID := chi.URLParam(r, "userID")
	intUserID, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(rw, "invalid user id ", http.StatusBadRequest)
		return
	}

	items, err := l.service.GetAll(int64(intUserID), key.Value)
	if err != nil {
		l.log.Sugar().Errorf("get all logo pass error: %v", err)
		http.Error(rw, apperrors.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(items)
}
