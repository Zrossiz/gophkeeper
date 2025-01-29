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
	GetAll(userID int64) ([]entities.Card, error)
}

func NewCardHandler(service CardService, logger *zap.Logger) *CardHandler {
	return &CardHandler{
		service: service,
		log:     logger,
	}
}

func (c *CardHandler) Create(rw http.ResponseWriter, r *http.Request) {
	var body dto.CreateCardDTO
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		fmt.Println(err)
		http.Error(rw, apperrors.ErrInvalidRequestBody, http.StatusBadRequest)
		return
	}

	err = c.service.Create(body)
	if err != nil {
		c.log.Sugar().Errorf("create card error: %v", err)
		http.Error(rw, apperrors.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func (c *CardHandler) Update(rw http.ResponseWriter, r *http.Request) {
	var body dto.UpdateCardDTO
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(rw, apperrors.ErrInvalidRequestBody, http.StatusBadRequest)
		return
	}

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

func (c *CardHandler) GetAll(rw http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	intUserID, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(rw, "invalid user id ", http.StatusBadRequest)
		return
	}

	items, err := c.service.GetAll(int64(intUserID))
	if err != nil {
		c.log.Sugar().Errorf("get all cards error: %v", err)
		http.Error(rw, apperrors.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(items)
}
