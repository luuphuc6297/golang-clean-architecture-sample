package repositories

import (
	"clean-architecture-api/internal/domain/entities"
	"context"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	BaseRepository[entities.User]
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
}
