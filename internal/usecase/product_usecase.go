package usecase

import (
	"clean-architecture-api/internal/domain/entities"
	"clean-architecture-api/internal/domain/repositories"
	"clean-architecture-api/pkg/logger"
	"context"

	"github.com/google/uuid"
)

// ProductUseCase defines the business logic interface for product operations
type ProductUseCase interface {
	Create(ctx context.Context, product *entities.Product, userID uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Product, error)
	Update(ctx context.Context, product *entities.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*entities.Product, error)
	GetByCategory(ctx context.Context, category string, limit, offset int) ([]*entities.Product, error)
}

type productUseCase struct {
	BaseUseCase
	productRepo repositories.ProductRepository
}

func NewProductUseCase(productRepo repositories.ProductRepository, logger logger.Logger) ProductUseCase {
	return &productUseCase{
		BaseUseCase: *NewBaseUseCase(logger),
		productRepo: productRepo,
	}
}

func (uc *productUseCase) Create(ctx context.Context, product *entities.Product, userID uuid.UUID) error {
	product.CreatedBy = userID

	if err := uc.productRepo.Create(ctx, product, userID); err != nil {
		return uc.HandleError(err, "failed to create product")
	}

	return nil
}

func (uc *productUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entities.Product, error) {
	// For GetByID operations, we need userID from context or parameter
	// For now, we'll extract from context or use a default approach
	userID := uc.getUserIDFromContext(ctx)

	product, err := uc.productRepo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, uc.HandleError(err, "product not found")
	}
	return product, nil
}

func (uc *productUseCase) Update(ctx context.Context, product *entities.Product) error {
	userID := uc.getUserIDFromContext(ctx)

	existingProduct, err := uc.productRepo.GetByID(ctx, product.ID, userID)
	if err != nil {
		return uc.HandleError(err, "product not found")
	}

	uc.updateProductFields(existingProduct, product)

	if err := uc.productRepo.Update(ctx, existingProduct, userID); err != nil {
		return uc.HandleError(err, "failed to update product")
	}

	return nil
}

func (uc *productUseCase) updateProductFields(existingProduct, product *entities.Product) {
	existingProduct.Name = product.Name
	existingProduct.Description = product.Description
	existingProduct.Price = product.Price
	existingProduct.Stock = product.Stock
	existingProduct.Category = product.Category
}

func (uc *productUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	userID := uc.getUserIDFromContext(ctx)

	if err := uc.ValidateEntityExists(ctx, func() error {
		_, err := uc.productRepo.GetByID(ctx, id, userID)
		return err
	}, "product"); err != nil {
		return err
	}

	if err := uc.productRepo.Delete(ctx, id, userID); err != nil {
		return uc.HandleError(err, "failed to delete product")
	}

	return nil
}

func (uc *productUseCase) List(ctx context.Context, limit, offset int) ([]*entities.Product, error) {
	userID := uc.getUserIDFromContext(ctx)

	products, err := uc.productRepo.List(ctx, limit, offset, userID)
	if err != nil {
		return nil, uc.HandleError(err, "failed to list products")
	}
	return products, nil
}

func (uc *productUseCase) GetByCategory(ctx context.Context, category string, limit, offset int) ([]*entities.Product,
	error,
) {
	products, err := uc.productRepo.GetByCategory(ctx, category, limit, offset)
	if err != nil {
		return nil, uc.HandleError(err, "failed to get products by category")
	}
	return products, nil
}

// getUserIDFromContext extracts user ID from context
func (uc *productUseCase) getUserIDFromContext(ctx context.Context) uuid.UUID {
	if userID, exists := ctx.Value("user_id").(uuid.UUID); exists {
		return userID
	}
	// Return a nil UUID if not found - this should be handled at middleware level
	return uuid.Nil
}
