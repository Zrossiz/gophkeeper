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
	GetAll(userID int64) ([]entities.LogoPassword, error)
}

func NewLogoPassHandler(service LogoPassService, logger *zap.Logger) *LogoPassHandler {
	return &LogoPassHandler{
		service: service,
		log:     logger,
	}
}

func (l *LogoPassHandler) Create(rw http.ResponseWriter, r *http.Request) {
	var body dto.CreateLogoPassDTO
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(rw, apperrors.ErrInvalidRequestBody, http.StatusBadRequest)
		return
	}

	err = l.service.Create(body)
	if err != nil {
		l.log.Sugar().Errorf("create logo pass error: %v", err)
		http.Error(rw, apperrors.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func (l *LogoPassHandler) Update(rw http.ResponseWriter, r *http.Request) {
	var body dto.UpdateLogoPassDTO
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(rw, apperrors.ErrInvalidRequestBody, http.StatusBadRequest)
		return
	}

	userID := chi.URLParam(r, "userID")
	intUserID, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(rw, "invalid user id ", http.StatusBadRequest)
		return
	}

	err = l.service.Update(int64(intUserID), body)
	if err != nil {
		l.log.Sugar().Errorf("update logo pass error: %v", err)
		http.Error(rw, apperrors.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (l *LogoPassHandler) GetAll(rw http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	intUserID, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(rw, "invalid user id ", http.StatusBadRequest)
		return
	}

	items, err := l.service.GetAll(int64(intUserID))
	if err != nil {
		l.log.Sugar().Errorf("get all logo pass error: %v", err)
		http.Error(rw, apperrors.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(items)
}
