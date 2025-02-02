package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Zrossiz/gophkeeper/internal/apperrors"
	"github.com/Zrossiz/gophkeeper/internal/dto"
	"github.com/Zrossiz/gophkeeper/internal/entities"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type CardHandler struct {
	service CardService
	log     *zap.Logger
}

type CardService interface {
	Create(body dto.CreateCardDTO) error
	Update(cardID int64, body dto.UpdateCardDTO) error
	GetAll(userID int64, key string) ([]entities.Card, error)
}

func NewCardHandler(service CardService, logger *zap.Logger) *CardHandler {
	return &CardHandler{
		service: service,
		log:     logger,
	}
}

// @Summary Создать карточку
// @Description Создает новую карточку пользователя
// @Tags card
// @Accept json
// @Produce json
// @Param body body dto.CreateCardDTO true "Данные для создания карточки"
// @Success 201 {string} string "Created"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /card [post]
// @Security BearerAuth
func (c *CardHandler) Create(rw http.ResponseWriter, r *http.Request) {
	key, err := r.Cookie("key")
	if err != nil {
		http.Error(rw, "key not found", http.StatusBadRequest)
		return
	}

	var body dto.CreateCardDTO
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		fmt.Println(err)
		http.Error(rw, apperrors.ErrInvalidRequestBody, http.StatusBadRequest)
		return
	}
	body.Key = key.Value

	err = c.service.Create(body)
	if err != nil {
		c.log.Sugar().Errorf("create card error: %v", err)
		http.Error(rw, apperrors.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

// @Summary Обновить карточку
// @Description Обновляет данные карточки
// @Tags card
// @Accept json
// @Produce json
// @Param cardID path int true "ID карточки"
// @Param body body dto.UpdateCardDTO true "Данные для обновления карточки"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /card/{cardID} [put]
// @Security BearerAuth
func (c *CardHandler) Update(rw http.ResponseWriter, r *http.Request) {
	key, err := r.Cookie("key")
	if err != nil {
		http.Error(rw, "key not found", http.StatusBadRequest)
		return
	}

	var body dto.UpdateCardDTO
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(rw, apperrors.ErrInvalidRequestBody, http.StatusBadRequest)
		return
	}
	body.Key = key.Value

	cardID := chi.URLParam(r, "cardID")
	intCardID, err := strconv.Atoi(cardID)
	if err != nil {
		http.Error(rw, "invalid user id ", http.StatusBadRequest)
		return
	}

	err = c.service.Update(int64(intCardID), body)
	if err != nil {
		c.log.Sugar().Errorf("update card error: %v", err)
		http.Error(rw, apperrors.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

// @Summary Получить все карточки пользователя
// @Description Возвращает список всех карточек пользователя
// @Tags card
// @Accept json
// @Produce json
// @Param userID path int true "ID пользователя"
// @Success 200 {array} entities.Card "Список карточек"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /card/user/{userID} [get]
// @Security BearerAuth
func (c *CardHandler) GetAll(rw http.ResponseWriter, r *http.Request) {
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

	items, err := c.service.GetAll(int64(intUserID), key.Value)
	if err != nil {
		c.log.Sugar().Errorf("get all cards error: %v", err)
		http.Error(rw, apperrors.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(items)
}
