package repository

import (
	"context"

	"clean-architecture-api/internal/domain/entities"
	"clean-architecture-api/internal/domain/repositories"
	"clean-architecture-api/pkg/logger"

	"gorm.io/gorm"
)

type productRepository struct {
	*CleanBaseRepositoryImpl[entities.Product]
}

// NewProductRepository creates a new product repository instance
func NewProductRepository(
	db *gorm.DB,
	authService repositories.AuthorizationService,
	auditLogger repositories.AuditLogger,
	logger logger.Logger,
) repositories.ProductRepository {
	return &productRepository{
		CleanBaseRepositoryImpl: NewCleanBaseRepository[entities.Product](db, auditLogger, logger, "product", authService),
	}
}

func (r *productRepository) GetByCategory(ctx context.Context, category string, limit, offset int) ([]*entities.Product, error) {
	var products []*entities.Product
	err := r.GetDB().WithContext(ctx).Where("category = ?", category).Limit(limit).Offset(offset).Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}
