// Package config provides functionality for loading and managing application configuration.
// It supports reading configuration values from environment variables with fallback defaults.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config represents the application configuration.
// It includes settings for server addresses, secrets, token durations, database URI, logging level, and cost.
type Config struct {
	BffAddress           string        // Address for the BFF (Backend for Frontend) server.
	ServerAddress        string        // Address for the main server.
	AccessSecret         string        // Secret key for access token generation.
	RefreshSecret        string        // Secret key for refresh token generation.
	DurationAccessToken  time.Duration // Duration for which access tokens are valid.
	DurationRefreshToken time.Duration // Duration for which refresh tokens are valid.
	DBURI                string        // URI for connecting to the database.
	LoggerLevel          string        // Logging level (e.g., DEBUG, INFO, ERROR).
	Cost                 int           // Cost factor for cryptographic operations (e.g., bcrypt).
}

// New initializes and returns a new Config instance by reading values from environment variables.
// If an environment variable is not set, it falls back to default values.
// Returns an error if required configurations (e.g., token durations) are invalid.
func New() (*Config, error) {
	cfg := Config{}

	// Load configuration values from environment variables or use defaults.
	cfg.BffAddress = getStringEnvOrDefault("BFF_ADDRESS", "localhost:9000")
	cfg.ServerAddress = getStringEnvOrDefault("SERVER_ADDRESS", "localhost:8080")
	cfg.DBURI = getStringEnvOrDefault("DB_URI", "host=localhost port=5432 user=postgres password=root dbname=gophkeeper sslmode=disable")
	cfg.AccessSecret = getStringEnvOrDefault("ACCESS_SECRET", "access")
	cfg.RefreshSecret = getStringEnvOrDefault("REFRESH_SECRET", "refresh")
	cfg.LoggerLevel = getStringEnvOrDefault("LOGGER_LEVEL", "DEBUG")
	cfg.Cost = getIntEnvOrDefault("COST", 3)

	// Parse token durations from environment variables.
	durationAccessSecret := getStringEnvOrDefault("DURATION_ACCESS", "24h")
	parsedDurationAccess, err := time.ParseDuration(durationAccessSecret)
	if err != nil {
		return nil, fmt.Errorf("invalid access token duration")
	}
	cfg.DurationAccessToken = parsedDurationAccess

	durationRefreshSecret := getStringEnvOrDefault("DURATION_REFRESH", "720h")
	parsedDurationRefresh, err := time.ParseDuration(durationRefreshSecret)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token duration")
	}
	cfg.DurationRefreshToken = parsedDurationRefresh

	return &cfg, nil
}

// getStringEnvOrDefault retrieves the value of an environment variable as a string.
// If the environment variable is not set or is empty, it returns the provided default value.
func getStringEnvOrDefault(envName string, defaultValue string) string {
	envValue := os.Getenv(envName)
	if envValue != "" {
		return envValue
	}

	return defaultValue
}

// getIntEnvOrDefault retrieves the value of an environment variable as an integer.
// If the environment variable is not set, is empty, or cannot be parsed as an integer,
// it returns the provided default value.
func getIntEnvOrDefault(envName string, defaultValue int) int {
	envValue := os.Getenv(envName)
	if envValue == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(envValue)
	if err != nil {
		fmt.Printf("error parsing %v, %v\n", envValue, err)
		return defaultValue
	}

	return intValue
}
