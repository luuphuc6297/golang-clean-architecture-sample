package repository

import (
	"clean-architecture-api/internal/domain/repositories"
	"clean-architecture-api/pkg/logger"
	"context"
	"errors"
	"fmt"

	domainerrors "clean-architecture-api/internal/domain/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CleanBaseRepositoryImpl[T any] struct {
	db           *gorm.DB
	auditLogger  repositories.AuditLogger
	logger       logger.Logger
	resourceName string
	authService  repositories.AuthorizationService
}

func NewCleanBaseRepository[T any](
	db *gorm.DB,
	auditLogger repositories.AuditLogger,
	logger logger.Logger,
	resourceName string,
	authService repositories.AuthorizationService,
) *CleanBaseRepositoryImpl[T] {
	return &CleanBaseRepositoryImpl[T]{
		db:           db,
		auditLogger:  auditLogger,
		logger:       logger,
		resourceName: resourceName,
		authService:  authService,
	}
}

func (r *CleanBaseRepositoryImpl[T]) Create(ctx context.Context, entity *T, userID uuid.UUID) error {
	if err := r.ValidateAccess(ctx, userID, "create"); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		r.logger.Error("Database create operation failed", err)
		return r.handleDatabaseError(err, "create", r.resourceName)
	}

	return r.AuditLog(ctx, userID, "create", entity)
}

func (r *CleanBaseRepositoryImpl[T]) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*T, error) {
	if err := r.ValidateAccess(ctx, userID, "read"); err != nil {
		return nil, err
	}

	var entity T
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.NewNotFoundError(
				fmt.Sprintf("%s_NOT_FOUND", r.resourceName),
				fmt.Sprintf("%s not found", r.resourceName),
			)
		}
		r.logger.Error("Database read operation failed", err)
		return nil, r.handleDatabaseError(err, "read", r.resourceName)
	}

	if err := r.AuditLog(ctx, userID, "read", &entity); err != nil {
		r.logger.Error("Failed to audit log read operation", err)
	}

	return &entity, nil
}

// Update updates an existing entity in the database
func (r *CleanBaseRepositoryImpl[T]) Update(ctx context.Context, entity *T, userID uuid.UUID) error {
	if err := r.ValidateAccess(ctx, userID, "update"); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Save(entity).Error; err != nil {
		r.logger.Error("Database update operation failed", err)
		return r.handleDatabaseError(err, "update", r.resourceName)
	}

	return r.AuditLog(ctx, userID, "update", entity)
}

func (r *CleanBaseRepositoryImpl[T]) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	if err := r.ValidateAccess(ctx, userID, "delete"); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Delete(new(T), "id = ?", id).Error; err != nil {
		r.logger.Error("Database delete operation failed", err)
		return r.handleDatabaseError(err, "delete", r.resourceName)
	}

	return r.AuditLog(ctx, userID, "delete", nil)
}

func (r *CleanBaseRepositoryImpl[T]) List(ctx context.Context, limit, offset int, userID uuid.UUID) ([]*T, error) {
	if err := r.ValidateAccess(ctx, userID, "list"); err != nil {
		return nil, err
	}

	var entities []*T
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&entities).Error
	if err != nil {
		r.logger.Error("Database list operation failed", err)
		return nil, r.handleDatabaseError(err, "list", r.resourceName)
	}

	if err := r.AuditLog(ctx, userID, "list", nil); err != nil {
		r.logger.Error("Failed to audit log list operation", err)
	}

	return entities, nil
}

func (r *CleanBaseRepositoryImpl[T]) ValidateAccess(ctx context.Context, userID uuid.UUID, action string) error {
	if r.authService == nil {
		return nil
	}
	return r.authService.CheckPermission(ctx, userID, r.resourceName, action)
}

func (r *CleanBaseRepositoryImpl[T]) AuditLog(ctx context.Context, userID uuid.UUID, action string, _ *T) error {
	if r.auditLogger == nil {
		return nil
	}

	resource := r.resourceName + ":" + action
	return r.auditLogger.LogAccess(ctx, userID, action, resource, uuid.Nil)
}

func (r *CleanBaseRepositoryImpl[T]) handleDatabaseError(err error, operation, resource string) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domainerrors.NewNotFoundError(
			fmt.Sprintf("%s_NOT_FOUND", resource),
			fmt.Sprintf("%s not found", resource),
		)
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return domainerrors.NewConflictError(
			fmt.Sprintf("%s_ALREADY_EXISTS", resource),
			fmt.Sprintf("%s already exists", resource),
		)
	}

	return domainerrors.NewDatabaseError(
		fmt.Sprintf("%s_%s_FAILED", operation, resource),
		fmt.Sprintf("database %s operation failed for %s", operation, resource),
		err,
	)
}

func (r *CleanBaseRepositoryImpl[T]) GetDB() *gorm.DB {
	return r.db
}
