package repository

import (
	"context"

	"clean-architecture-api/internal/domain/entities"
	"clean-architecture-api/internal/domain/repositories"
	"clean-architecture-api/pkg/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type policySQLiteRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewPolicySQLiteRepository creates a new SQLite policy repository instance
func NewPolicySQLiteRepository(db *gorm.DB, logger logger.Logger) repositories.PolicyRepository {
	return &policySQLiteRepository{
		db:     db,
		logger: logger,
	}
}

func (r *policySQLiteRepository) Create(ctx context.Context, policy *entities.PolicyDocument) error {
	policySQLite := entities.FromPolicyDocument(policy)

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		policyToCreate := *policySQLite
		policyToCreate.Statements = nil

		if err := tx.Create(&policyToCreate).Error; err != nil {
			return err
		}

		for _, stmt := range policySQLite.Statements {
			stmt.ID = uuid.New().String()
			stmt.PolicyID = policyToCreate.ID
			if err := tx.Create(&stmt).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *policySQLiteRepository) GetByRole(ctx context.Context, role string) ([]*entities.PolicyDocument, error) {
	var policiesSQLite []*entities.PolicyDocumentSQLite

	err := r.db.WithContext(ctx).
		Preload("Statements").
		Joins("JOIN policy_statements ON policy_documents.id = policy_statements.policy_id").
		Where("policy_statements.principal = ? OR policy_statements.principal = ?", "role:"+role, "*").
		Where("policy_documents.is_active = ?", true).
		Find(&policiesSQLite).Error

	if err != nil {
		return nil, err
	}

	policies := make([]*entities.PolicyDocument, len(policiesSQLite))
	for i, policySQLite := range policiesSQLite {
		policies[i] = policySQLite.ToPolicyDocument()
	}

	return r.deduplicatePolicies(policies), nil
}

func (r *policySQLiteRepository) GetActive(ctx context.Context) ([]*entities.PolicyDocument, error) {
	var policiesSQLite []*entities.PolicyDocumentSQLite

	err := r.db.WithContext(ctx).
		Preload("Statements").
		Where("is_active = ?", true).
		Find(&policiesSQLite).Error

	if err != nil {
		return nil, err
	}

	// Convert to regular entities
	policies := make([]*entities.PolicyDocument, len(policiesSQLite))
	for i, policySQLite := range policiesSQLite {
		policies[i] = policySQLite.ToPolicyDocument()
	}

	return policies, nil
}

func (r *policySQLiteRepository) Update(ctx context.Context, policy *entities.PolicyDocument) error {
	policySQLite := entities.FromPolicyDocument(policy)

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(policySQLite).Error; err != nil {
			return err
		}

		if err := tx.Where("policy_id = ?", policySQLite.ID).Delete(&entities.PolicyStatementSQLite{}).Error; err != nil {
			return err
		}

		for i := range policySQLite.Statements {
			policySQLite.Statements[i].ID = uuid.New().String()
			policySQLite.Statements[i].PolicyID = policySQLite.ID
		}

		if len(policySQLite.Statements) > 0 {
			if err := tx.Create(&policySQLite.Statements).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *policySQLiteRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("policy_id = ?", id.String()).Delete(&entities.PolicyStatementSQLite{}).Error; err != nil {
			return err
		}

		if err := tx.Delete(&entities.PolicyDocumentSQLite{}, "id = ?", id.String()).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *policySQLiteRepository) deduplicatePolicies(policies []*entities.PolicyDocument) []*entities.PolicyDocument {
	seen := make(map[uuid.UUID]bool)
	var result []*entities.PolicyDocument

	for _, policy := range policies {
		if !seen[policy.ID] {
			seen[policy.ID] = true
			result = append(result, policy)
		}
	}

	return result
}
