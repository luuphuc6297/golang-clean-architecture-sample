package database

import (
	"clean-architecture-api/internal/domain/constants"
	"clean-architecture-api/internal/domain/entities"
	"clean-architecture-api/internal/infrastructure/auth"
	"clean-architecture-api/pkg/logger"
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func NewDatabase() (*gorm.DB, error) {
	config, err := NewDatabaseConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load database config: %w", err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Ho_Chi_Minh",
		config.Host, config.User, config.Password, config.Name, config.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entities.User{},
		&entities.Product{},
		&entities.PolicyDocument{},
		&entities.PolicyStatement{},
		&auth.AuditLogEntry{},
	)
}

func InitializeDefaultPolicies(db *gorm.DB, logger logger.Logger) error {
	return initializePoliciesWithModel(db, logger, &entities.PolicyDocument{}, func() error {
		ctx := context.Background()
		policies := []*entities.PolicyDocument{
			createAdminPolicy(),
			createUserPolicy(),
		}

		for _, policy := range policies {
			if err := createPolicyWithStatements(ctx, db, policy); err != nil {
				logger.Error("Failed to create policy: "+policy.Name, err)
				return err
			}
			logger.Info("Created policy: " + policy.Name)
		}
		return nil
	})
}

func createAdminPolicy() *entities.PolicyDocument {
	return &entities.PolicyDocument{
		ID:       uuid.New(),
		Name:     "admin-full-access",
		Version:  "1.0",
		IsActive: true,
		Statements: []entities.PolicyStatement{
			{
				ID:         uuid.New(),
				Effect:     constants.PolicyEffectAllow,
				Principal:  "role:" + constants.RoleAdmin,
				Action:     "*",
				Resource:   "*",
				Conditions: map[string]interface{}{},
			},
		},
	}
}

func createUserPolicy() *entities.PolicyDocument {
	statements := []entities.PolicyStatement{}

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

	for i, permission := range productPermissions {
		statements = append(statements, entities.PolicyStatement{
			ID:         uuid.New(),
			Effect:     constants.PolicyEffectAllow,
			Principal:  "role:" + constants.RoleUser,
			Action:     userActions[i],
			Resource:   permission,
			Conditions: map[string]interface{}{},
		})
	}

	return &entities.PolicyDocument{
		ID:         uuid.New(),
		Name:       "user-product-access",
		Version:    "1.0",
		IsActive:   true,
		Statements: statements,
	}
}

func createPolicyWithStatements(ctx context.Context, db *gorm.DB, policy *entities.PolicyDocument) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(policy).Error; err != nil {
			return err
		}

		for i := range policy.Statements {
			policy.Statements[i].ID = uuid.New()
			policy.Statements[i].PolicyID = policy.ID
		}

		if len(policy.Statements) > 0 {
			if err := tx.Create(&policy.Statements).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func initializePoliciesWithModel(db *gorm.DB, logger logger.Logger, model interface{}, createFunc func() error) error {
	var count int64
	if err := db.Model(model).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		logger.Info("Policies already exist, skipping initialization")
		return nil
	}

	if err := createFunc(); err != nil {
		return err
	}

	logger.Info("Default policies initialized successfully")
	return nil
}
