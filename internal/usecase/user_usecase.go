package usecase

import (
	"context"

	"clean-architecture-api/internal/domain/entities"
	domainerrors "clean-architecture-api/internal/domain/errors"
	"clean-architecture-api/internal/domain/repositories"
	"clean-architecture-api/pkg/logger"

	"github.com/google/uuid"
)

type UserUseCase interface {
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entities.User, error)
	Update(ctx context.Context, user *entities.User, userID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	List(ctx context.Context, limit, offset int, userID uuid.UUID) ([]*entities.User, error)
}

type userUseCase struct {
	BaseUseCase
	userRepo repositories.UserRepository
}

func NewUserUseCase(userRepo repositories.UserRepository, logger logger.Logger) UserUseCase {
	return &userUseCase{
		BaseUseCase: *NewBaseUseCase(logger),
		userRepo:    userRepo,
	}
}

func (uc *userUseCase) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entities.User, error) {
	user, err := uc.userRepo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, domainerrors.ErrUserNotFound
	}
	return user, nil
}

func (uc *userUseCase) Update(ctx context.Context, user *entities.User, userID uuid.UUID) error {
	existingUser, err := uc.userRepo.GetByID(ctx, user.ID, userID)
	if err != nil {
		return domainerrors.ErrUserNotFound
	}

	uc.updateUserFields(existingUser, user)

	if err := uc.userRepo.Update(ctx, existingUser, userID); err != nil {
		return uc.HandleError(err, "failed to update user")
	}

	return nil
}

func (uc *userUseCase) updateUserFields(existingUser, user *entities.User) {
	existingUser.FirstName = user.FirstName
	existingUser.LastName = user.LastName
	existingUser.Role = user.Role
	existingUser.IsActive = user.IsActive
}

func (uc *userUseCase) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	if err := uc.ValidateEntityExists(ctx, func() error {
		_, err := uc.userRepo.GetByID(ctx, id, userID)
		return err
	}, "user"); err != nil {
		return err
	}

	if err := uc.userRepo.Delete(ctx, id, userID); err != nil {
		return domainerrors.ErrDeleteUser
	}

	return nil
}

func (uc *userUseCase) List(ctx context.Context, limit, offset int, userID uuid.UUID) ([]*entities.User, error) {
	users, err := uc.userRepo.List(ctx, limit, offset, userID)
	if err != nil {
		return nil, uc.HandleError(err, "failed to list users")
	}
	return users, nil
}
