package postgres

import (
	"testing"

	"github.com/Zrossiz/gophkeeper/internal/dto"
	"github.com/Zrossiz/gophkeeper/internal/entities"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestNotesStorage_Create(t *testing.T) {
	db, cleanup := setupDB(t)
	defer cleanup()

	storage := NewNotesStorage(db)

	// Создаём заметку
	body := dto.CreateNoteDTO{
		UserID:   1,
		Title:    "Test Title",
		TextData: "Test Note Data",
	}

	// Сохраняем заметку
	err := storage.Create(body)
	assert.NoError(t, err, "Create should insert a note without error")

	// Проверяем, что данные были сохранены в базе
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM notes WHERE user_id = $1 AND title = $2", body.UserID, body.Title).Scan(&count)
	assert.NoError(t, err, "Failed to query notes table")
	assert.Equal(t, 1, count, "Expected one note to be inserted")
}

func TestNotesStorage_Update(t *testing.T) {
	db, cleanup := setupDB(t)
	defer cleanup()

	storage := NewNotesStorage(db)

	// Сначала создаём заметку для обновления
	createBody := dto.CreateNoteDTO{
		UserID:   1,
		Title:    "Test Title",
		TextData: "Test Note Data",
	}
	err := storage.Create(createBody)
	assert.NoError(t, err, "Create should insert a note without error")

	// Обновляем заметку
	updateBody := dto.UpdateNoteDTO{
		Title:    "Updated Title",
		TextData: "Updated Note Data",
	}
	err = storage.Update(1, updateBody) // Предположим, что ID заметки - 1
	assert.NoError(t, err, "Update should update the note without error")

	// Проверяем, что данные были обновлены
	var updatedNote entities.Note
	err = db.QueryRow("SELECT title, text_data FROM notes WHERE id = $1", 1).
		Scan(&updatedNote.Title, &updatedNote.TextData)
	assert.NoError(t, err, "Failed to query updated note")

	// Проверяем обновлённые данные
	assert.Equal(t, updateBody.Title, updatedNote.Title, "Title should be updated")
	assert.Equal(t, updateBody.TextData, updatedNote.TextData, "TextData should be updated")
}

func TestNotesStorage_GetAllByUser(t *testing.T) {
	db, cleanup := setupDB(t)
	defer cleanup()

	storage := NewNotesStorage(db)

	// Создаём несколько заметок для одного пользователя
	userID := int64(1)
	notesDTOs := []dto.CreateNoteDTO{
		{UserID: int(userID), Title: "Title 1", TextData: "Note data 1"},
		{UserID: int(userID), Title: "Title 2", TextData: "Note data 2"},
	}

	// Сохраняем все заметки
	for _, body := range notesDTOs {
		err := storage.Create(body)
		assert.NoError(t, err, "Create should insert a note without error")
	}

	// Получаем все заметки для пользователя
	notes, err := storage.GetAllByUser(int(userID))
	assert.NoError(t, err, "GetAllByUser should not return an error")
	assert.Len(t, notes, len(notesDTOs), "Expected the same number of notes")

	// Проверяем, что данные сохранены корректно
	for i, note := range notes {
		assert.Equal(t, notesDTOs[i].UserID, note.UserID, "UserID should match")
		assert.Equal(t, notesDTOs[i].Title, note.Title, "Title should match")
		assert.Equal(t, notesDTOs[i].TextData, note.TextData, "TextData should match")
	}
}
