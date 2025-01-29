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

type BinaryHandler struct {
	service BinaryService
	log     *zap.Logger
}

type BinaryService interface {
	Create(body dto.CreateBinaryDTO) error
	Update(userID int, body dto.UpdateBinaryDTO) error
	GetAll(userID int64) ([]entities.BinaryData, error)
}

func NewBinaryHandler(service BinaryService, logger *zap.Logger) *BinaryHandler {
	return &BinaryHandler{
		service: service,
		log:     logger,
	}
}

func (b *BinaryHandler) Create(rw http.ResponseWriter, r *http.Request) {
	var body dto.CreateBinaryDTO
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(rw, apperrors.ErrInvalidRequestBody, http.StatusBadRequest)
		return
	}

	err = b.service.Create(body)
	if err != nil {
		b.log.Sugar().Errorf("create binary error: %v", err)
		http.Error(rw, apperrors.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func (b *BinaryHandler) Update(rw http.ResponseWriter, r *http.Request) {
	var body dto.UpdateBinaryDTO
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(rw, apperrors.ErrInvalidRequestBody, http.StatusBadRequest)
		return
	}

	binaryID := chi.URLParam(r, "binaryID")
	intBinaryID, err := strconv.Atoi(binaryID)
	if err != nil {
		http.Error(rw, "invalid user id ", http.StatusBadRequest)
		return
	}

	err = b.service.Update(intBinaryID, body)
	if err != nil {
		b.log.Sugar().Errorf("update binary error: %v", err)
		http.Error(rw, apperrors.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (b *BinaryHandler) GetAll(rw http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	intUserID, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(rw, "invalid user id ", http.StatusBadRequest)
		return
	}

	items, err := b.service.GetAll(int64(intUserID))
	if err != nil {
		b.log.Sugar().Errorf("error get all binaries: %v", err)
		http.Error(rw, apperrors.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(items)
}
