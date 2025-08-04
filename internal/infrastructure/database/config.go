// Package database provides database configuration and connection management.
// It supports both PostgreSQL and SQLite databases with environment-based configuration.
package database

import (
	"clean-architecture-api/internal/domain/constants"
	"fmt"
	"os"
)

// DatabaseConfig holds configuration for database connection
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// NewDatabaseConfig creates a new database configuration from environment variables
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

// SQLiteConfig holds configuration for SQLite database connection
type SQLiteConfig struct {
	DBPath string
}

// NewSQLiteConfig creates a new SQLite configuration from environment variables
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
