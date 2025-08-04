package repositories

import (
	"clean-architecture-api/internal/domain/entities"
	"context"

	"github.com/google/uuid"
)

type AuthorizationService interface {
	CheckPermission(ctx context.Context, userID uuid.UUID, resource, action string) error
	CheckResourcePermission(ctx context.Context, userID uuid.UUID, resource, action, resourceID string) error
	GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]entities.Permission, error)
	GetEffectivePermissions(ctx context.Context, userID uuid.UUID) ([]entities.Permission, error)
	QuickCheck(userRole, resource, action string) bool
	ValidateRole(userRole string) error
	GetAllowedActionsForRole(userRole, resource string) ([]string, error)
	CreateEnrichedContext(ctx context.Context, userID uuid.UUID, role, email string) context.Context
}

type AuditLogger interface {
	LogAccess(ctx context.Context, userID uuid.UUID, action, resource string, entityID uuid.UUID) error
	LogDataAccess(ctx context.Context, userID uuid.UUID, action, resource string, data interface{}) error
}

type BaseRepository[T any] interface {
	Create(ctx context.Context, entity *T, userID uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*T, error)
	Update(ctx context.Context, entity *T, userID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	List(ctx context.Context, limit, offset int, userID uuid.UUID) ([]*T, error)

	ValidateAccess(ctx context.Context, userID uuid.UUID, action string) error
	AuditLog(ctx context.Context, userID uuid.UUID, action string, entity *T) error
}
