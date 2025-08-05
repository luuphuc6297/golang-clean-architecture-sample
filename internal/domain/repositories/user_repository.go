package repositories

import (
	"clean-architecture-api/internal/domain/entities"
	"context"
)

type UserRepository interface {
	BaseRepository[entities.User]
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
}
