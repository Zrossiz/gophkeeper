package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/Zrossiz/gophkeeper/internal/dto"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupDB(t *testing.T) (*sql.DB, func()) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpassword",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err, "Failed to start PostgreSQL container")
	t.Cleanup(func() {
		err := container.Terminate(ctx)
		require.NoError(t, err, "Failed to stop container")
	})

	host, err := container.Host(ctx)
	require.NoError(t, err, "Failed to get container host")

	port, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err, "Failed to get mapped port")

	dsn := fmt.Sprintf("postgres://testuser:testpassword@%s:%s/testdb?sslmode=disable", host, port.Port())

	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err, "Failed to connect to PostgreSQL")

	path := "/Users/zrossiz/Desktop/GoProjects/praktikum/projects/gophkeeper/migrations/1_init.sql"

	query, _ := os.ReadFile(path)

	_, err = db.Exec(string(query))
	require.NoError(t, err, "Failed to migrate")

	insertUserQuery := `INSERT INTO users (username, password) VALUES ('test', 'test')`
	_, err = db.Exec(insertUserQuery)
	require.NoError(t, err, "Failed to create user")

	return db, func() {
		db.Close()
	}
}

func TestBinaryStorage_Create(t *testing.T) {
	db, cleanup := setupDB(t)
	defer cleanup()

	storage := NewBinaryStorage(db)

	body := dto.SetStorageBinaryDTO{
		UserID: 1,
		Title:  "test title",
		Data:   []byte("test data"),
	}

	err := storage.Create(context.Background(), body)
	assert.NoError(t, err, "Create should insert binary data without error")

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM binary_data WHERE user_id = $1 AND title = $2", body.UserID, body.Title).Scan(&count)
	assert.NoError(t, err, "Failed to query binary_data table")
	assert.Equal(t, 1, count, "Expected one row to be inserted")
}

func TestBinaryStorage_Update(t *testing.T) {
	db, cleanup := setupDB(t)
	defer cleanup()

	storage := NewBinaryStorage(db)

	createBody := dto.SetStorageBinaryDTO{
		UserID: 1,
		Title:  "test title",
		Data:   []byte("initial data"),
	}
	err := storage.Create(context.Background(), createBody)
	assert.NoError(t, err, "Create should insert binary data without error")

	updateBody := dto.SetStorageBinaryDTO{
		UserID: 1,
		Title:  "test title",
		Data:   []byte("updated data"),
	}
	err = storage.Update(context.Background(), updateBody)
	assert.NoError(t, err, "Update should update binary data without error")

	var data []byte
	err = db.QueryRow("SELECT binary_data FROM binary_data WHERE user_id = $1 AND title = $2", updateBody.UserID, updateBody.Title).Scan(&data)
	assert.NoError(t, err, "Failed to query binary_data table")
	assert.Equal(t, []byte("updated data"), data, "Expected data to be updated")
}

func TestBinaryStorage_GetAllByUser(t *testing.T) {
	db, cleanup := setupDB(t)
	defer cleanup()

	storage := NewBinaryStorage(db)

	userID := int64(1)
	bodies := []dto.SetStorageBinaryDTO{
		{UserID: int(userID), Title: "title1", Data: []byte("data1")},
		{UserID: int(userID), Title: "title2", Data: []byte("data2")},
	}
	for _, body := range bodies {
		err := storage.Create(context.Background(), body)
		assert.NoError(t, err, "Create should insert binary data without error")
	}

	binaryDataList, err := storage.GetAllByUser(context.Background(), userID)
	assert.NoError(t, err, "GetAllByUser should retrieve binary data without error")
	assert.Len(t, binaryDataList, len(bodies), "Expected number of binary data records to match")

	for i, binaryData := range binaryDataList {
		assert.Equal(t, bodies[i].UserID, binaryData.UserID, "UserID should match")
		assert.Equal(t, bodies[i].Title, binaryData.Title, "Title should match")
		assert.Equal(t, bodies[i].Data, binaryData.Data, "Data should match")
	}
}
