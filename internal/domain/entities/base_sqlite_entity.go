package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseSQLiteEntity provides common fields for SQLite entities with integer primary key
type BaseSQLiteEntity struct {
	ID        string         `json:"id" gorm:"type:text;primary_key"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// BeforeCreate is a GORM hook that runs before SQLite entity creation
func (e *BaseSQLiteEntity) BeforeCreate(_ *gorm.DB) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return nil
}
