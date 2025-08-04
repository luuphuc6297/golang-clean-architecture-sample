package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// PolicyStatementSQLite represents a policy statement entity for SQLite database
type PolicyStatementSQLite struct {
	ID         string    `json:"id" gorm:"type:text;primary_key"`
	PolicyID   string    `json:"policy_id" gorm:"type:text;not null"`
	Effect     string    `json:"effect" gorm:"not null"`
	Principal  string    `json:"principal" gorm:"not null"`
	Action     string    `json:"action" gorm:"not null"`
	Resource   string    `json:"resource" gorm:"not null"`
	Conditions string    `json:"conditions,omitempty" gorm:"type:text"` // JSON as string for SQLite
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// TableName returns the table name for PolicyStatementSQLite entity
func (PolicyStatementSQLite) TableName() string {
	return "policy_statements"
}

// PolicyDocumentSQLite represents a policy document entity for SQLite database
type PolicyDocumentSQLite struct {
	BaseSQLiteEntity
	Name       string                  `json:"name" gorm:"not null;unique"`
	Version    string                  `json:"version" gorm:"not null;default:1.0"`
	Statements []PolicyStatementSQLite `json:"statements" gorm:"foreignKey:PolicyID"`
	IsActive   bool                    `json:"is_active" gorm:"default:true"`
}

// TableName returns the table name for PolicyDocumentSQLite entity
func (PolicyDocumentSQLite) TableName() string {
	return "policy_documents"
}

// ToPolicyDocument converts SQLite policy document to domain policy document
func (p *PolicyDocumentSQLite) ToPolicyDocument() *PolicyDocument {
	id, _ := uuid.Parse(p.ID)

	statements := make([]PolicyStatement, len(p.Statements))
	for i, stmt := range p.Statements {
		stmtID, _ := uuid.Parse(stmt.ID)
		policyID, _ := uuid.Parse(stmt.PolicyID)

		var conditions map[string]interface{}
		if stmt.Conditions != "" {
			if err := json.Unmarshal([]byte(stmt.Conditions), &conditions); err != nil {
				conditions = map[string]interface{}{}
			}
		}

		statements[i] = PolicyStatement{
			ID:         stmtID,
			PolicyID:   policyID,
			Effect:     stmt.Effect,
			Principal:  stmt.Principal,
			Action:     stmt.Action,
			Resource:   stmt.Resource,
			Conditions: conditions,
			CreatedAt:  stmt.CreatedAt,
			UpdatedAt:  stmt.UpdatedAt,
		}
	}

	return &PolicyDocument{
		ID:         id,
		Name:       p.Name,
		Version:    p.Version,
		Statements: statements,
		IsActive:   p.IsActive,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}
}

// FromPolicyDocument converts domain policy document to SQLite policy document
func FromPolicyDocument(policy *PolicyDocument) *PolicyDocumentSQLite {
	statements := make([]PolicyStatementSQLite, len(policy.Statements))
	for i, stmt := range policy.Statements {
		conditionsJSON, _ := json.Marshal(stmt.Conditions)

		statements[i] = PolicyStatementSQLite{
			ID:         stmt.ID.String(),
			PolicyID:   policy.ID.String(),
			Effect:     stmt.Effect,
			Principal:  stmt.Principal,
			Action:     stmt.Action,
			Resource:   stmt.Resource,
			Conditions: string(conditionsJSON),
			CreatedAt:  stmt.CreatedAt,
			UpdatedAt:  stmt.UpdatedAt,
		}
	}

	return &PolicyDocumentSQLite{
		BaseSQLiteEntity: BaseSQLiteEntity{
			ID:        policy.ID.String(),
			CreatedAt: policy.CreatedAt,
			UpdatedAt: policy.UpdatedAt,
		},
		Name:       policy.Name,
		Version:    policy.Version,
		Statements: statements,
		IsActive:   policy.IsActive,
	}
}
