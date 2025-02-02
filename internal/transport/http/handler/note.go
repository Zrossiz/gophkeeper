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

type NoteHandler struct {
	service NoteService
	log     *zap.Logger
}

type NoteService interface {
	Create(body dto.CreateNoteDTO) error
	Update(noteID int, body dto.UpdateNoteDTO) error
	GetAll(userID int, key string) ([]entities.Note, error)
}

func NewNoteHandler(service NoteService, log *zap.Logger) *NoteHandler {
	return &NoteHandler{
		service: service,
		log:     log,
	}
}

// @Summary Создать заметку
// @Description Создает новую заметку
// @Tags note
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer токен" default(Bearer {token})
// @Param body body dto.CreateNoteDTO true "Данные для создания заметки"
// @Success 201 {string} string "Created"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /note [post]
// @Security BearerAuth
func (n *NoteHandler) Create(rw http.ResponseWriter, r *http.Request) {
	key, err := r.Cookie("key")
	if err != nil {
		http.Error(rw, "key not found", http.StatusBadRequest)
		return
	}

	var body dto.CreateNoteDTO
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(rw, apperrors.ErrInvalidRequestBody, http.StatusBadRequest)
		return
	}

	body.Key = key.Value

	err = n.service.Create(body)
	if err != nil {
		http.Error(rw, apperrors.ErrInternalServer, http.StatusInternalServerError)
		n.log.Sugar().Errorf("create note error: %v", err)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

// @Summary Обновить заметку
// @Description Обновляет существующую заметку
// @Tags note
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer токен" default(Bearer {token})
// @Param noteID path int true "ID заметки"
// @Param body body dto.UpdateNoteDTO true "Данные для обновления"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /note/{noteID} [put]
// @Security BearerAuth
func (n *NoteHandler) Update(rw http.ResponseWriter, r *http.Request) {
	noteID := chi.URLParam(r, "noteID")
	intNoteID, err := strconv.Atoi(noteID)
	if err != nil {
		http.Error(rw, "invalid note id ", http.StatusBadRequest)
		return
	}

	key, err := r.Cookie("key")
	if err != nil {
		http.Error(rw, "key not found", http.StatusBadRequest)
		return
	}

	var body dto.UpdateNoteDTO
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(rw, apperrors.ErrInvalidRequestBody, http.StatusBadRequest)
		return
	}

	body.Key = key.Value

	err = n.service.Update(intNoteID, body)
	if err != nil {
		n.log.Sugar().Errorf("update note id error: %v", err)
		http.Error(rw, apperrors.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

// @Summary Получить все заметки пользователя
// @Description Возвращает список всех заметок пользователя
// @Tags note
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer токен" default(Bearer {token})
// @Param userID path int true "ID пользователя"
// @Success 200 {array} entities.Note "Список заметок"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /note/user/{userID} [get]
// @Security BearerAuth
func (n *NoteHandler) GetAll(rw http.ResponseWriter, r *http.Request) {
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

	n.log.Info("url valid")

	items, err := n.service.GetAll(intUserID, key.Value)
	if err != nil {
		n.log.Sugar().Errorf("get all notes error: %v", err)
		http.Error(rw, apperrors.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(items)
}
