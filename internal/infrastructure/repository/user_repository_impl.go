package repository

import (
	"context"

	"clean-architecture-api/internal/domain/entities"
	"clean-architecture-api/internal/domain/repositories"
	"clean-architecture-api/pkg/logger"

	"gorm.io/gorm"
)

type userRepository struct {
	*CleanBaseRepositoryImpl[entities.User]
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(
	db *gorm.DB,
	authService repositories.AuthorizationService,
	auditLogger repositories.AuditLogger,
	logger logger.Logger,
) repositories.UserRepository {
	return &userRepository{
		CleanBaseRepositoryImpl: NewCleanBaseRepository[entities.User](db, auditLogger, logger, "user", authService),
	}
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	err := r.GetDB().WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
