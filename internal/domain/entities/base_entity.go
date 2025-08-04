// BaseEntity is a base struct for all entities providing common fields.

// Package entities contains domain entity definitions for the clean architecture API.
package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseEntity provides common fields for all entities including ID, timestamps and soft delete
type BaseEntity struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// BeforeCreate is a GORM hook that runs before entity creation to set the ID
func (e *BaseEntity) BeforeCreate(_ *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}
