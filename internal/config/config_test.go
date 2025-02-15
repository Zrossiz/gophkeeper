package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfigWithDefaultValues(t *testing.T) {
	t.Setenv("BFF_ADDRESS", "")
	t.Setenv("SERVER_ADDRESS", "")
	t.Setenv("DB_URI", "")
	t.Setenv("ACCESS_SECRET", "")
	t.Setenv("REFRESH_SECRET", "")
	t.Setenv("LOGGER_LEVEL", "")
	t.Setenv("COST", "")
	t.Setenv("DURATION_ACCESS", "")
	t.Setenv("DURATION_REFRESH", "")

	cfg, err := New()
	require.NoError(t, err, "Should create config without error")

	assert.Equal(t, "localhost:9000", cfg.BffAddress)
	assert.Equal(t, "localhost:8080", cfg.ServerAddress)
	assert.Equal(t, "host=localhost port=5432 user=postgres password=root dbname=gophkeeper sslmode=disable", cfg.DBURI)
	assert.Equal(t, "access", cfg.AccessSecret)
	assert.Equal(t, "refresh", cfg.RefreshSecret)
	assert.Equal(t, "DEBUG", cfg.LoggerLevel)
	assert.Equal(t, 3, cfg.Cost)

	expectedDurationAccess := 24 * time.Hour
	assert.Equal(t, expectedDurationAccess, cfg.DurationAccessToken)

	expectedDurationRefresh := 720 * time.Hour
	assert.Equal(t, expectedDurationRefresh, cfg.DurationRefreshToken)
}

func TestNewConfigWithEnvValues(t *testing.T) {
	t.Setenv("BFF_ADDRESS", "localhost:9090")
	t.Setenv("SERVER_ADDRESS", "localhost:8081")
	t.Setenv("DB_URI", "host=localhost port=5432 user=test password=test dbname=testdb sslmode=disable")
	t.Setenv("ACCESS_SECRET", "customAccessSecret")
	t.Setenv("REFRESH_SECRET", "customRefreshSecret")
	t.Setenv("LOGGER_LEVEL", "INFO")
	t.Setenv("COST", "10")
	t.Setenv("DURATION_ACCESS", "48h")
	t.Setenv("DURATION_REFRESH", "1000h")

	cfg, err := New()
	require.NoError(t, err, "Should create config without error")

	assert.Equal(t, "localhost:9090", cfg.BffAddress)
	assert.Equal(t, "localhost:8081", cfg.ServerAddress)
	assert.Equal(t, "host=localhost port=5432 user=test password=test dbname=testdb sslmode=disable", cfg.DBURI)
	assert.Equal(t, "customAccessSecret", cfg.AccessSecret)
	assert.Equal(t, "customRefreshSecret", cfg.RefreshSecret)
	assert.Equal(t, "INFO", cfg.LoggerLevel)
	assert.Equal(t, 10, cfg.Cost)

	expectedDurationAccess := 48 * time.Hour
	assert.Equal(t, expectedDurationAccess, cfg.DurationAccessToken)

	expectedDurationRefresh := 1000 * time.Hour
	assert.Equal(t, expectedDurationRefresh, cfg.DurationRefreshToken)
}

func TestNewConfigWithInvalidDuration(t *testing.T) {
	t.Setenv("DURATION_ACCESS", "invalidDuration")
	t.Setenv("DURATION_REFRESH", "invalidDuration")

	_, err := New()
	require.Error(t, err, "Should return an error for invalid duration")
	assert.Contains(t, err.Error(), "invalid access token duration")
}

func TestGetStringEnvOrDefault(t *testing.T) {
	assert.Equal(t, "localhost:9000", getStringEnvOrDefault("BFF_ADDRESS", "localhost:9000"))
	assert.Equal(t, "localhost:8080", getStringEnvOrDefault("SERVER_ADDRESS", "localhost:8080"))

	t.Setenv("BFF_ADDRESS", "localhost:9999")
	assert.Equal(t, "localhost:9999", getStringEnvOrDefault("BFF_ADDRESS", "localhost:9000"))
}

func TestGetIntEnvOrDefault(t *testing.T) {
	assert.Equal(t, 3, getIntEnvOrDefault("COST", 3))

	t.Setenv("COST", "5")
	assert.Equal(t, 5, getIntEnvOrDefault("COST", 3))

	t.Setenv("COST", "invalid")
	assert.Equal(t, 3, getIntEnvOrDefault("COST", 3))
}
