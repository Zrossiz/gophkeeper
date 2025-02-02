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
	GetAll(userID int64, key string) ([]entities.BinaryData, error)
}

func NewBinaryHandler(service BinaryService, logger *zap.Logger) *BinaryHandler {
	return &BinaryHandler{
		service: service,
		log:     logger,
	}
}

// @Summary Загрузить бинарные данные
// @Description Загружает бинарный файл пользователя
// @Tags binary
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Bearer токен" default(Bearer {token})
// @Param file formData file true "Файл для загрузки"
// @Param user_id formData int true "ID пользователя"
// @Success 201 {string} string "File uploaded successfully!"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /binary/ [post]
// @Security BearerAuth
func (b *BinaryHandler) Create(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(rw, "File too large or invalid request", http.StatusBadRequest)
		return
	}

	key, err := r.Cookie("key")
	if err != nil {
		http.Error(rw, "key not found", http.StatusBadRequest)
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

	userID := r.FormValue("user_id")
	intUserID, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(rw, "invalid user id", http.StatusBadRequest)
	}

	body := dto.CreateBinaryDTO{
		Title:  header.Filename,
		Data:   fileData,
		UserID: intUserID,
		Key:    key.Value,
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

// @Summary Получить все бинарные данные пользователя
// @Description Возвращает список всех загруженных бинарных данных пользователя
// @Tags binary
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer токен" default(Bearer {token})
// @Param userID path int true "ID пользователя"
// @Success 200 {array} entities.BinaryData "Список бинарных данных"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /binary/user/{userID} [get]
// @Security BearerAuth
func (b *BinaryHandler) GetAll(rw http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	intUserID, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(rw, "invalid user id ", http.StatusBadRequest)
		return
	}

	key, err := r.Cookie("key")
	if err != nil {
		http.Error(rw, "key not found", http.StatusBadRequest)
		return
	}

	items, err := b.service.GetAll(int64(intUserID), key.Value)
	if err != nil {
		b.log.Sugar().Errorf("error get all binaries: %v", err)
		http.Error(rw, apperrors.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(items)
}
