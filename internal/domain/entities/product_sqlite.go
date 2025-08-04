package entities

import (
	"clean-architecture-api/internal/domain/constants"
	"clean-architecture-api/internal/domain/validators"

	"github.com/google/uuid"
)

type ProductSQLite struct {
	BaseSQLiteEntity
	Name        string  `json:"name" gorm:"not null"`
	Description string  `json:"description"`
	Price       float64 `json:"price" gorm:"not null"`
	Stock       int     `json:"stock" gorm:"default:0"`
	Category    string  `json:"category"`
	CreatedBy   string  `json:"created_by" gorm:"type:text"`
}

func (ProductSQLite) TableName() string {
	return "products"
}

func (p *ProductSQLite) Validate() error {
	if err := validators.ValidateRequired(constants.FieldName, p.Name); err != nil {
		return err
	}
	if err := validators.ValidatePrice(p.Price); err != nil {
		return err
	}
	if err := validators.ValidateStock(p.Stock); err != nil {
		return err
	}
	return nil
}

func (p *ProductSQLite) ToProduct() *Product {
	id, _ := uuid.Parse(p.ID)
	createdBy, _ := uuid.Parse(p.CreatedBy)
	product := &Product{
		BaseEntity: BaseEntity{
			ID:        id,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
			DeletedAt: p.DeletedAt,
		},
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		Category:    p.Category,
		CreatedBy:   createdBy,
	}
	return product
}

func FromProduct(product *Product) *ProductSQLite {
	return &ProductSQLite{
		BaseSQLiteEntity: BaseSQLiteEntity{
			ID:        product.ID.String(),
			CreatedAt: product.CreatedAt,
			UpdatedAt: product.UpdatedAt,
			DeletedAt: product.DeletedAt,
		},
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Category:    product.Category,
		CreatedBy:   product.CreatedBy.String(),
	}
}
