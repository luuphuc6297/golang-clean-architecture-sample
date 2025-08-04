package repositories

import (
	"context"

	"clean-architecture-api/internal/domain/entities"
)

type UserRepository interface {
	BaseRepository[entities.User]
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
}
