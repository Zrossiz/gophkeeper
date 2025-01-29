package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	BffAddress           string
	ServerAddress        string
	AccessSecret         string
	RefreshSecret        string
	DurationAccessToken  time.Duration
	DurationRefreshToken time.Duration
	DBURI                string
	LoggerLevel          string
	Cost                 int
}

func New() (*Config, error) {
	cfg := Config{}

	cfg.BffAddress = getStringEnvOrDefault("BFF_ADDRESS", "localhost:9000")
	cfg.ServerAddress = getStringEnvOrDefault("SERVER_ADDRESS", "localhost:8080")
	cfg.DBURI = getStringEnvOrDefault("DB_URI", "host=localhost port=5432 user=postgres password=root dbname=gophkeerper sslmode=disable")
	cfg.AccessSecret = getStringEnvOrDefault("ACCESS_SECRET", "access")
	cfg.RefreshSecret = getStringEnvOrDefault("REFRESH_SECRET", "refresh")
	cfg.LoggerLevel = getStringEnvOrDefault("LOGGER_LEVEL", "DEBUG")
	cfg.Cost = getIntEnvOrDefault("COST", 3)

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

func getStringEnvOrDefault(envName string, defaultValue string) string {
	envValue := os.Getenv(envName)
	if envValue != "" {
		return envValue
	}

	return defaultValue
}

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
