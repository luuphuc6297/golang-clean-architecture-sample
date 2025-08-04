package repositories

import (
	"context"

	"clean-architecture-api/internal/domain/entities"
)

type ProductRepository interface {
	BaseRepository[entities.Product]
	GetByCategory(ctx context.Context, category string, limit, offset int) ([]*entities.Product, error)
}
