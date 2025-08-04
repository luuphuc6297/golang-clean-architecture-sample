package database

import (
	"log"

	"clean-architecture-api/internal/domain/entities"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewInMemoryDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&entities.UserSQLite{},
		&entities.ProductSQLite{},
		&entities.PolicyDocumentSQLite{},
		&entities.PolicyStatementSQLite{},
	); err != nil {
		return nil, err
	}

	log.Println("In-memory database initialized and migrated successfully")
	return db, nil
}
