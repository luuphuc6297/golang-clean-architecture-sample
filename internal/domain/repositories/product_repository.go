package repositories

import (
	"clean-architecture-api/internal/domain/entities"
	"context"
)

// ProductRepository defines the interface for product data operations
type ProductRepository interface {
	BaseRepository[entities.Product]
	GetByCategory(ctx context.Context, category string, limit, offset int) ([]*entities.Product, error)
}
