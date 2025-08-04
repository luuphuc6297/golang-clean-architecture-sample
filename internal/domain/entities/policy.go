package entities

import (
	"clean-architecture-api/internal/domain/constants"
	"time"

	"github.com/google/uuid"
)

// PolicyStatement represents a single policy statement with effect, actions, and resources
type PolicyStatement struct {
	ID         uuid.UUID              `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	PolicyID   uuid.UUID              `json:"policy_id" gorm:"type:uuid;not null"`
	Effect     string                 `json:"effect" gorm:"not null"`
	Principal  string                 `json:"principal" gorm:"not null"`
	Action     string                 `json:"action" gorm:"not null"`
	Resource   string                 `json:"resource" gorm:"not null"`
	Conditions map[string]interface{} `json:"conditions,omitempty" gorm:"type:jsonb"`
	CreatedAt  time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time              `json:"updated_at" gorm:"autoUpdateTime"`
}

// PolicyDocument represents a complete policy document containing multiple statements
type PolicyDocument struct {
	ID         uuid.UUID         `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name       string            `json:"name" gorm:"not null;unique"`
	Version    string            `json:"version" gorm:"not null;default:'1.0'"`
	Statements []PolicyStatement `json:"statements" gorm:"foreignKey:PolicyID"`
	IsActive   bool              `json:"is_active" gorm:"default:true"`
	CreatedAt  time.Time         `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time         `json:"updated_at" gorm:"autoUpdateTime"`
}

// PermissionRequest represents a request for permission evaluation
type PermissionRequest struct {
	UserID     uuid.UUID              `json:"user_id"`
	Role       string                 `json:"role"`
	Resource   string                 `json:"resource"`
	Action     string                 `json:"action"`
	ResourceID string                 `json:"resource_id,omitempty"`
	Context    map[string]interface{} `json:"context"`
}

// PermissionResponse represents the result of a permission evaluation
type PermissionResponse struct {
	Allowed  bool                   `json:"allowed"`
	Reason   string                 `json:"reason,omitempty"`
	Policies []string               `json:"policies,omitempty"`
	Context  map[string]interface{} `json:"context,omitempty"`
}

// Permission represents a permission with resource and action
type Permission struct {
	Resource   string `json:"resource"`
	Action     string `json:"action"`
	Role       string `json:"role"`
	ResourceID string `json:"resource_id,omitempty"`
}

// IsValid validates that the policy statement has required fields
func (ps *PolicyStatement) IsValid() bool {
	return ps.Effect == constants.PolicyEffectAllow || ps.Effect == constants.PolicyEffectDeny
}
