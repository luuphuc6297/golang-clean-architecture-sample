// Package main provides the entry point for the Clean Architecture API server.
package main

import (
	"clean-architecture-api/internal/delivery/http"
	"clean-architecture-api/internal/infrastructure/database"
	"clean-architecture-api/pkg/logger"
	"clean-architecture-api/pkg/newrelic"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	logger := logger.NewLogger()

	if err := loadEnv(); err != nil {
		logger.Fatal("Failed to load environment variables", err)
	}

	nrConfig := newrelic.NewConfig()
	nrApp, err := newrelic.NewApplication(nrConfig)
	if err != nil {
		logger.Warn("Failed to initialize New Relic application", err)
	} else if nrApp != nil {
		logger.Info("New Relic application initialized successfully")
	}
	db, err := database.NewDatabaseWithNewRelic(nrApp)
	if err != nil {
		logger.Fatal("Failed to connect to database", err)
	}
	if err := database.InitializeDefaultPolicies(db, logger); err != nil {
		logger.Fatal("Failed to initialize default policies", err)
	}
	server, err := http.NewServerWithNewRelic(db, logger, nrApp)
	if err != nil {
		logger.Fatal("Failed to create HTTP server", err)
	}
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
	envFile := ".env"
	if _, err := os.Stat(envFile); err == nil {
		if err := godotenv.Load(envFile); err != nil {
			return err
		}
	}
	if os.Getenv("ENV") == "" {
		if err := os.Setenv("ENV", "development"); err != nil {
			return err
		}
	}
	return nil
}
