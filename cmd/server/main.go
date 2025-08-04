package main

import (
	"os"

	"clean-architecture-api/internal/delivery/http"
	"clean-architecture-api/internal/infrastructure/database"
	"clean-architecture-api/pkg/logger"
)

func main() {
	// Initialize logger
	logger := logger.NewLogger()

	// Load environment variables
	if err := loadEnv(); err != nil {
		logger.Fatal("Failed to load environment variables", err)
	}

	// Initialize database
	db, err := database.NewDatabase()
	if err != nil {
		logger.Fatal("Failed to connect to database", err)
	}

	// Initialize default policies
	if err := database.InitializeDefaultPolicies(db, logger); err != nil {
		logger.Fatal("Failed to initialize default policies", err)
	}

	// Initialize HTTP server
	server, err := http.NewServer(db, logger)
	if err != nil {
		logger.Fatal("Failed to create HTTP server", err)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("Server starting on port " + port)
	if err := server.Run(":" + port); err != nil {
		logger.Fatal("Failed to start server", err)
	}
}

func loadEnv() error {
	// In production, you might want to use a proper config management library
	// For now, we'll use godotenv for development
	if os.Getenv("ENV") == "" {
		os.Setenv("ENV", "development")
	}
	return nil
}
