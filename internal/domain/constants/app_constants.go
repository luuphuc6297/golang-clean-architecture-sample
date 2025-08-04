// Package constants defines application-wide constant values used throughout the system.
package constants

const (
	// DefaultLimit represents the default pagination limit
	DefaultLimit  = 10
	// DefaultOffset represents the default pagination offset
	DefaultOffset = 0
	// MaxLimit represents the maximum allowed pagination limit
	MaxLimit      = 100

	// RoleUser represents the standard user role
	RoleUser  = "user"
	// RoleAdmin represents the administrator role
	RoleAdmin = "admin"

	// JWTAccessTokenDuration represents access token duration in minutes
	JWTAccessTokenDuration  = 15
	// JWTRefreshTokenDuration represents refresh token duration in days
	JWTRefreshTokenDuration = 7

	// DefaultDBHost represents the default database host
	DefaultDBHost = "localhost"
	// DefaultDBPort represents the default database port
	DefaultDBPort = "5432"
	// DefaultDBUser represents the default database user
	DefaultDBUser = "postgres"
	// DefaultDBName represents the default database name
	DefaultDBName = "clean_architecture_api"

	// DefaultPort represents the default server port
	DefaultPort = "8080"
	// DefaultEnv represents the default environment
	DefaultEnv  = "development"

	// SystemUserID represents the system user ID for internal operations
	SystemUserID = "00000000-0000-0000-0000-000000000000"
)
