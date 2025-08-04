//go:build sqlite
// +build sqlite

package main

import (
	"log"
	"os"

	"clean-architecture-api/internal/delivery/http"
	"clean-architecture-api/internal/infrastructure/database"
	"clean-architecture-api/pkg/logger"
)

func main() {
	logger := logger.NewLogger()

	if err := loadEnv(); err != nil {
		logger.Fatal("Failed to load environment variables", err)
	}

	db, err := database.NewSQLiteDatabase()
	if err != nil {
		logger.Fatal("Failed to connect to SQLite database", err)
	}

	// Initialize default policies for SQLite
	if err := database.InitializeSQLiteDefaultPolicies(db, logger); err != nil {
		logger.Fatal("Failed to initialize default policies", err)
	}

	server, err := http.NewServer(db, logger)
	if err != nil {
		logger.Fatal("Failed to create server", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("Server starting on port " + port + " with SQLite database")
	if err := server.Run(":" + port); err != nil {
		log.Fatal("Failed to start server", err)
	}
}

func loadEnv() error {
	if os.Getenv("ENV") == "" {
		os.Setenv("ENV", "development")
	}

	/** Remove this after testing */
	if os.Getenv("JWT_SECRET_KEY") == "" {
		os.Setenv("JWT_SECRET_KEY", "dev-secret-key-change-in-production")
	}

	// Set default port if not provided
	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "8080")
	}

	// Set default SQLite database path if not provided
	if os.Getenv("SQLITE_DB_PATH") == "" {
		os.Setenv("SQLITE_DB_PATH", "./data/clean_architecture_api.db")
	}

	// Set default log level if not provided
	if os.Getenv("LOG_LEVEL") == "" {
		os.Setenv("LOG_LEVEL", "info")
	}

	return nil
}
