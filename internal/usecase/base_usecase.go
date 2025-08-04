package usecase

import (
	"context"
	"fmt"

	domainerrors "clean-architecture-api/internal/domain/errors"
	"clean-architecture-api/pkg/logger"
)

type BaseUseCase struct {
	logger logger.Logger
}

func NewBaseUseCase(logger logger.Logger) *BaseUseCase {
	return &BaseUseCase{logger: logger}
}

func (uc *BaseUseCase) HandleError(err error, message string) error {
	uc.logger.Error(message, err)
	return fmt.Errorf("%s: %w", message, err)
}

func (uc *BaseUseCase) HandleDatabaseError(err error, operation, entity string) error {
	uc.logger.Error(fmt.Sprintf("Database operation failed: %s %s", operation, entity), err)
	return domainerrors.NewDatabaseError(
		fmt.Sprintf("%s_%s_FAILED", operation, entity),
		fmt.Sprintf("failed to %s %s", operation, entity),
		err,
	)
}

func (uc *BaseUseCase) HandleNotFoundError(entity string) error {
	message := fmt.Sprintf("%s not found", entity)
	code := fmt.Sprintf("%s_NOT_FOUND", entity)
	return domainerrors.NewNotFoundError(code, message)
}

func (uc *BaseUseCase) ValidateEntityExists(_ context.Context, getFunc func() error, entityName string) error {
	if err := getFunc(); err != nil {
		return uc.HandleError(err, entityName+" not found")
	}
	return nil
}
