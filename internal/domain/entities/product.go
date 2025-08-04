package entities

import (
	"clean-architecture-api/internal/domain/constants"
	"clean-architecture-api/internal/domain/validators"

	"github.com/google/uuid"
)

// Product represents a product entity in the system
type Product struct {
	BaseEntity
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Price       float64   `json:"price" gorm:"not null"`
	Stock       int       `json:"stock" gorm:"default:0"`
	Category    string    `json:"category"`
	CreatedBy   uuid.UUID `json:"created_by" gorm:"type:uuid"`
}

// TableName returns the database table name for the Product entity
func (Product) TableName() string {
	return "products"
}

// Validate validates the product entity fields and returns an error if invalid
func (p *Product) Validate() error {
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
