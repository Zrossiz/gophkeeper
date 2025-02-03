package postgres

import (
	"testing"

	"github.com/Zrossiz/gophkeeper/internal/dto"
	"github.com/Zrossiz/gophkeeper/internal/entities"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestLogoPassStorage_CreateLogoPass(t *testing.T) {
	db, cleanup := setupDB(t)
	defer cleanup()

	storage := NewLogoPassStorage(db)

	body := dto.CreateLogoPassDTO{
		UserId:   1,
		AppName:  "TestApp",
		Username: "testuser",
		Password: "testpassword",
	}

	err := storage.CreateLogoPass(body)
	assert.NoError(t, err, "CreateLogoPass should not return an error")

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM passwords WHERE user_id = $1 AND app_name = $2", body.UserId, body.AppName).Scan(&count)
	assert.NoError(t, err, "Failed to query passwords table")
	assert.Equal(t, 1, count, "Expected one logo pass to be inserted")
}

func TestLogoPassStorage_GetAllByUser(t *testing.T) {
	db, cleanup := setupDB(t)
	defer cleanup()

	storage := NewLogoPassStorage(db)

	userID := int64(1)
	logoPassDTOs := []dto.CreateLogoPassDTO{
		{UserId: int(userID), AppName: "App1", Username: "user1", Password: "pass1"},
		{UserId: int(userID), AppName: "App2", Username: "user2", Password: "pass2"},
	}

	for _, body := range logoPassDTOs {
		err := storage.CreateLogoPass(body)
		assert.NoError(t, err, "CreateLogoPass should not return an error")
	}

	logoPasses, err := storage.GetAllByUser(userID)
	assert.NoError(t, err, "GetAllByUser should not return an error")
	assert.Len(t, logoPasses, len(logoPassDTOs), "Expected the same number of logo passes")

	for i, logoPass := range logoPasses {
		assert.Equal(t, logoPassDTOs[i].UserId, logoPass.UserID, "UserID should match")
		assert.Equal(t, logoPassDTOs[i].AppName, logoPass.AppName, "AppName should match")
		assert.Equal(t, logoPassDTOs[i].Username, logoPass.Username, "Username should match")
		assert.Equal(t, logoPassDTOs[i].Password, logoPass.Password, "Password should match")
	}
}

func TestLogoPassStorage_UpdateLogoPass(t *testing.T) {
	db, cleanup := setupDB(t)
	defer cleanup()

	storage := NewLogoPassStorage(db)

	body := dto.CreateLogoPassDTO{
		UserId:   1,
		AppName:  "TestApp",
		Username: "testuser",
		Password: "testpassword",
	}

	err := storage.CreateLogoPass(body)
	assert.NoError(t, err, "CreateLogoPass should not return an error")

	updateBody := dto.UpdateLogoPassDTO{
		Username: "updateduser",
		Password: "updatedpassword",
	}

	err = storage.UpdateLogoPass(1, updateBody) // Предположим, что ID логина - 1
	assert.NoError(t, err, "UpdateLogoPass should not return an error")

	var updatedLogoPass entities.LogoPassword
	err = db.QueryRow("SELECT username, password FROM passwords WHERE id = $1", 1).
		Scan(&updatedLogoPass.Username, &updatedLogoPass.Password)
	assert.NoError(t, err, "Failed to query updated logo pass")

	assert.Equal(t, updateBody.Username, updatedLogoPass.Username, "Username should be updated")
	assert.Equal(t, updateBody.Password, updatedLogoPass.Password, "Password should be updated")
}
