// Package constants defines application-wide constant values used throughout the system.
package constants

const (
	DefaultLimit  = 10
	DefaultOffset = 0
	MaxLimit      = 100

	RoleUser  = "user"
	RoleAdmin = "admin"

	JWTAccessTokenDuration  = 15
	JWTRefreshTokenDuration = 7

	DefaultDBHost = "localhost"
	DefaultDBPort = "5432"
	DefaultDBUser = "postgres"
	DefaultDBName = "clean_architecture_api"

	DefaultPort = "8080"
	DefaultEnv  = "development"

	SystemUserID = "00000000-0000-0000-0000-000000000000"
)
