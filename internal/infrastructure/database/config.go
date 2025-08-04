package database

import (
	"clean-architecture-api/internal/domain/constants"
	"fmt"
	"os"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func NewDatabaseConfig() (*DatabaseConfig, error) {
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("DB_PASSWORD environment variable is required")
	}

	return &DatabaseConfig{
		Host:     getEnvOrDefault("DB_HOST", constants.DefaultDBHost),
		Port:     getEnvOrDefault("DB_PORT", constants.DefaultDBPort),
		User:     getEnvOrDefault("DB_USER", constants.DefaultDBUser),
		Password: password,
		Name:     getEnvOrDefault("DB_NAME", constants.DefaultDBName),
	}, nil
}

type SQLiteConfig struct {
	DBPath string
}

func NewSQLiteConfig() *SQLiteConfig {
	return &SQLiteConfig{
		DBPath: getEnvOrDefault("SQLITE_DB_PATH", "./data/clean_architecture_api.db"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
