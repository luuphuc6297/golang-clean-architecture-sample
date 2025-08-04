package database

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"clean-architecture-api/internal/domain/constants"
	"clean-architecture-api/internal/domain/entities"
	"clean-architecture-api/pkg/logger"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func NewSQLiteDatabase() (*gorm.DB, error) {
	config := NewSQLiteConfig()

	if err := os.MkdirAll("./data", 0o755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	db, err := gorm.Open(sqlite.Open(config.DBPath), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQLite database: %w", err)
	}

	if err := autoMigrateSQLite(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

func autoMigrateSQLite(db *gorm.DB) error {
	return db.AutoMigrate(
		&entities.UserSQLite{},
		&entities.ProductSQLite{},
		&entities.PolicyDocumentSQLite{},
		&entities.PolicyStatementSQLite{},
	)
}

func InitializeSQLiteDefaultPolicies(db *gorm.DB, logger logger.Logger) error {
	return initializePoliciesWithModel(db, logger, &entities.PolicyDocumentSQLite{}, func() error {
		ctx := context.Background()
		policies := []*entities.PolicyDocumentSQLite{
			createSQLiteAdminPolicy(),
			createSQLiteUserPolicy(),
		}

		for _, policy := range policies {
			if err := createSQLitePolicyWithStatements(ctx, db, policy); err != nil {
				logger.Error("Failed to create policy: "+policy.Name, err)
				return err
			}
			logger.Info("Created policy: " + policy.Name)
		}
		return nil
	})
}

func createSQLiteAdminPolicy() *entities.PolicyDocumentSQLite {
	conditions, _ := json.Marshal(map[string]interface{}{})

	return &entities.PolicyDocumentSQLite{
		BaseSQLiteEntity: entities.BaseSQLiteEntity{
			ID: uuid.New().String(),
		},
		Name:     "admin-full-access",
		Version:  "1.0",
		IsActive: true,
		Statements: []entities.PolicyStatementSQLite{
			{
				Effect:     constants.PolicyEffectAllow,
				Principal:  "role:" + constants.RoleAdmin,
				Action:     "*",
				Resource:   "*",
				Conditions: string(conditions),
			},
		},
	}
}

func createSQLiteUserPolicy() *entities.PolicyDocumentSQLite {
	statements := []entities.PolicyStatementSQLite{}

	productPermissions := []string{
		constants.PermissionProductCreate,
		constants.PermissionProductRead,
		constants.PermissionProductUpdate,
		constants.PermissionProductDelete,
		constants.PermissionProductList,
	}

	userActions := []string{
		constants.ActionCreate,
		constants.ActionRead,
		constants.ActionUpdate,
		constants.ActionDelete,
		constants.ActionList,
	}

	conditions, _ := json.Marshal(map[string]interface{}{})

	for i, permission := range productPermissions {
		statements = append(statements, entities.PolicyStatementSQLite{
			Effect:     constants.PolicyEffectAllow,
			Principal:  "role:" + constants.RoleUser,
			Action:     userActions[i],
			Resource:   permission,
			Conditions: string(conditions),
		})
	}

	return &entities.PolicyDocumentSQLite{
		BaseSQLiteEntity: entities.BaseSQLiteEntity{
			ID: uuid.New().String(),
		},
		Name:       "user-product-access",
		Version:    "1.0",
		IsActive:   true,
		Statements: statements,
	}
}

func createSQLitePolicyWithStatements(ctx context.Context, db *gorm.DB, policy *entities.PolicyDocumentSQLite) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		policyToCreate := *policy
		policyToCreate.Statements = nil

		if err := tx.Create(&policyToCreate).Error; err != nil {
			return err
		}

		for _, stmt := range policy.Statements {
			stmt.ID = uuid.New().String()
			stmt.PolicyID = policyToCreate.ID
			if err := tx.Create(&stmt).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
