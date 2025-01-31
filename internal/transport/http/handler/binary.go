package handler

import (
	"encoding/json"
	"fmt"
	"io"
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
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(rw, "File too large or invalid request", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(rw, "File is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		http.Error(rw, "Failed to read file", http.StatusInternalServerError)
		return
	}

	body := dto.CreateBinaryDTO{
		Title: header.Filename,
		Data:  fileData,
	}

	err = b.service.Create(body)
	if err != nil {
		b.log.Sugar().Errorf("create binary error: %v", err)
		http.Error(rw, apperrors.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	fmt.Fprintln(rw, "File uploaded successfully!")
}

func (b *BinaryHandler) Update(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(rw, "File too large or invalid request", http.StatusBadRequest)
		return
	}

	binaryID := chi.URLParam(r, "binaryID")
	intBinaryID, err := strconv.Atoi(binaryID)
	if err != nil {
		http.Error(rw, "Invalid binary ID", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	var fileData []byte
	if err == nil {
		defer file.Close()
		fileData, err = io.ReadAll(file)
		if err != nil {
			http.Error(rw, "Failed to read file", http.StatusInternalServerError)
			return
		}
	}

	body := dto.UpdateBinaryDTO{
		Title: header.Filename,
		Data:  fileData,
	}

	err = b.service.Update(intBinaryID, body)
	if err != nil {
		b.log.Sugar().Errorf("update binary error: %v", err)
		http.Error(rw, apperrors.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	fmt.Fprintln(rw, "File updated successfully!")
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
