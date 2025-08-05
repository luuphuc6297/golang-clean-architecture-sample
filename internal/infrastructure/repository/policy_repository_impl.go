package repository

import (
	"clean-architecture-api/internal/domain/entities"
	"clean-architecture-api/internal/domain/repositories"
	"clean-architecture-api/pkg/logger"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type policyRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewPolicyRepository(db *gorm.DB, logger logger.Logger) repositories.PolicyRepository {
	return &policyRepository{
		db:     db,
		logger: logger,
	}
}

func (r *policyRepository) Create(ctx context.Context, policy *entities.PolicyDocument) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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

func (r *policyRepository) GetByRole(ctx context.Context, role string) ([]*entities.PolicyDocument, error) {
	var policies []*entities.PolicyDocument

	err := r.db.WithContext(ctx).
		Preload("Statements").
		Joins("JOIN policy_statements ON policy_documents.id = policy_statements.policy_id").
		Where("policy_statements.principal = ? OR policy_statements.principal = ?", "role:"+role, "*").
		Where("policy_documents.is_active = ?", true).
		Find(&policies).Error
	if err != nil {
		return nil, err
	}

	return r.deduplicatePolicies(policies), nil
}

func (r *policyRepository) GetActive(ctx context.Context) ([]*entities.PolicyDocument, error) {
	var policies []*entities.PolicyDocument

	err := r.db.WithContext(ctx).
		Preload("Statements").
		Where("is_active = ?", true).
		Find(&policies).Error
	if err != nil {
		return nil, err
	}

	return policies, nil
}

func (r *policyRepository) Update(ctx context.Context, policy *entities.PolicyDocument) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(policy).Error; err != nil {
			return err
		}

		if err := tx.Where("policy_id = ?", policy.ID).Delete(&entities.PolicyStatement{}).Error; err != nil {
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

func (r *policyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("policy_id = ?", id).Delete(&entities.PolicyStatement{}).Error; err != nil {
			return err
		}

		if err := tx.Delete(&entities.PolicyDocument{}, "id = ?", id).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *policyRepository) deduplicatePolicies(policies []*entities.PolicyDocument) []*entities.PolicyDocument {
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
