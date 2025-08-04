package repositories

import (
	"clean-architecture-api/internal/domain/entities"
	"context"

	"github.com/google/uuid"
)

// PolicyEngine defines the interface for policy evaluation and management
type PolicyEngine interface {
	Evaluate(ctx context.Context, req *entities.PermissionRequest) (*entities.PermissionResponse, error)
	LoadPolicies(ctx context.Context) error
	AddPolicy(ctx context.Context, policy *entities.PolicyDocument) error
	RemovePolicy(ctx context.Context, policyID uuid.UUID) error
	GetPoliciesForRole(ctx context.Context, role string) ([]*entities.PolicyDocument, error)
}

// PolicyRepository defines the interface for policy data operations
type PolicyRepository interface {
	Create(ctx context.Context, policy *entities.PolicyDocument) error
	GetByRole(ctx context.Context, role string) ([]*entities.PolicyDocument, error)
	GetActive(ctx context.Context) ([]*entities.PolicyDocument, error)
	Update(ctx context.Context, policy *entities.PolicyDocument) error
	Delete(ctx context.Context, id uuid.UUID) error
}
