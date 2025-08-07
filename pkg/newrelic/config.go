// Package newrelic provides New Relic application monitoring integration.
package newrelic

import (
	"fmt"
	"os"

	"github.com/newrelic/go-agent/v3/newrelic"
)

// Config holds New Relic configuration options.
type Config struct {
	AppName    string
	LicenseKey string
	Enabled    bool
}

// NewConfig creates a new New Relic configuration from environment variables.
func NewConfig() *Config {
	enabled := os.Getenv("NEW_RELIC_ENABLED")
	if enabled != "true" {
		return &Config{Enabled: false}
	}

	return &Config{
		AppName:    getEnvOrDefault("NEW_RELIC_APP_NAME", "clean-architecture-api"),
		LicenseKey: os.Getenv("NEW_RELIC_LICENSE_KEY"),
		Enabled:    true,
	}
}

// NewApplication creates a new New Relic application instance.
func NewApplication(cfg *Config) (*newrelic.Application, error) {
	if !cfg.Enabled {
		return nil, nil
	}

	if cfg.LicenseKey == "" {
		return nil, fmt.Errorf("NEW_RELIC_LICENSE_KEY is required when New Relic is enabled")
	}

	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(cfg.AppName),
		newrelic.ConfigLicense(cfg.LicenseKey),
		newrelic.ConfigAppLogForwardingEnabled(true),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create New Relic application: %w", err)
	}

	return app, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
