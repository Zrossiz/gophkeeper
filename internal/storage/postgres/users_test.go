package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/Zrossiz/gophkeeper/internal/dto"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var db *sql.DB
var cleanup func()

func TestMain(m *testing.M) {

	var err error
	db, cleanup, err = setupUserDB()
	if err != nil {
		fmt.Printf("Error setting up DB: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()

	cleanup()

	os.Exit(code)
}

func setupUserDB() (*sql.DB, func(), error) {
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
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start PostgreSQL container: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get mapped port: %w", err)
	}

	dsn := fmt.Sprintf("postgres://testuser:testpassword@%s:%s/testdb?sslmode=disable", host, port.Port())

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	migrationPath := "/Users/zrossiz/Desktop/GoProjects/praktikum/projects/gophkeeper/migrations/1_init.sql"

	query, err := os.ReadFile(migrationPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read migration file: %w", err)
	}

	_, err = db.Exec(string(query))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to execute migration: %w", err)
	}

	return db, func() {
		container.Terminate(ctx)
		db.Close()
	}, nil
}

func clearDB() {
	_, err := db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE;")
	if err != nil {
		fmt.Printf("Error clearing DB: %v\n", err)
	}
}

func TestUserStorage_Create(t *testing.T) {
	storage := NewUserStorage(db)

	user := dto.UserDTO{
		Username: "testuser",
		Password: "testpassword",
	}

	err := storage.Create(context.Background(), user)
	assert.NoError(t, err, "Create should insert a user without error")

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", user.Username).Scan(&count)
	assert.NoError(t, err, "Failed to query users table")
	assert.Equal(t, 1, count, "Expected one user to be inserted")

	clearDB()
}

func TestUserStorage_GetUserByUsername(t *testing.T) {
	storage := NewUserStorage(db)

	createUser := dto.UserDTO{
		Username: "testuser",
		Password: "testpassword",
	}
	err := storage.Create(context.Background(), createUser)
	assert.NoError(t, err, "Create should insert a user without error")

	user, err := storage.GetUserByUsername(context.Background(), "testuser")
	assert.NoError(t, err, "GetUserByUsername should return user without error")

	assert.Equal(t, createUser.Username, user.Username, "Username should match")
	assert.Equal(t, createUser.Password, user.Password, "Password should match")

	clearDB()
}
