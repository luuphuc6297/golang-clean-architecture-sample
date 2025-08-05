package repositories

import (
	"clean-architecture-api/internal/domain/entities"
	"context"
)

type ProductRepository interface {
	BaseRepository[entities.Product]
	GetByCategory(ctx context.Context, category string, limit, offset int) ([]*entities.Product, error)
}
