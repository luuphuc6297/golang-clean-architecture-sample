// Package usecase contains business logic implementations following clean architecture principles.
// Use cases orchestrate data flow between entities and repositories while enforcing business rules.
package usecase

import (
	"context"

	"clean-architecture-api/internal/domain/constants"
	"clean-architecture-api/internal/domain/entities"
	domainerrors "clean-architecture-api/internal/domain/errors"
	"clean-architecture-api/internal/domain/repositories"
	"clean-architecture-api/internal/domain/validators"
	"clean-architecture-api/internal/infrastructure/auth"
	"clean-architecture-api/pkg/logger"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase interface {
	Register(ctx context.Context, email, password, firstName, lastName string) (*entities.User, error)
	Login(ctx context.Context, email, password string) (*auth.TokenPair, error)
	RefreshToken(ctx context.Context, refreshToken string) (*auth.TokenPair, error)
	ValidateToken(ctx context.Context, token string) (*auth.Claims, error)
}

type authUseCase struct {
	BaseUseCase
	userRepo    repositories.UserRepository
	authService auth.AuthService
}

func NewAuthUseCase(userRepo repositories.UserRepository, authService auth.AuthService, logger logger.Logger) AuthUseCase {
	return &authUseCase{
		BaseUseCase: *NewBaseUseCase(logger),
		userRepo:    userRepo,
		authService: authService,
	}
}

func (uc *authUseCase) Register(ctx context.Context, email, password, firstName, lastName string) (*entities.User, error) {
	if err := validators.ValidateRegisterRequest(email, password, firstName, lastName); err != nil {
		uc.logger.Error("User registration failed: validation error", err.Error())
		return nil, err
	}

	if err := uc.checkUserExists(ctx, email); err != nil {
		return nil, err
	}

	hashedPassword, err := uc.hashPassword(password)
	if err != nil {
		return nil, err
	}

	user := uc.createUser(email, hashedPassword, firstName, lastName)

	systemUserID := uuid.MustParse(constants.SystemUserID)
	if err := uc.userRepo.Create(ctx, user, systemUserID); err != nil {
		uc.logger.Error("Failed to create user in database", err.Error())
		return nil, domainerrors.ErrFailedToCreateUser
	}

	uc.logger.Info("User registered successfully", email)
	return user, nil
}

func (uc *authUseCase) checkUserExists(ctx context.Context, email string) error {
	existingUser, err := uc.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		uc.logger.Error("User registration failed: user already exists", email)
		return domainerrors.ErrUserAlreadyExists
	}
	return nil
}

func (uc *authUseCase) hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", domainerrors.ErrFailedToProcessPassword
	}
	return string(hashedPassword), nil
}

func (uc *authUseCase) createUser(email, hashedPassword, firstName, lastName string) *entities.User {
	return &entities.User{
		Email:     email,
		Password:  hashedPassword,
		FirstName: firstName,
		LastName:  lastName,
		Role:      "user",
		IsActive:  true,
	}
}

func (uc *authUseCase) Login(ctx context.Context, email, password string) (*auth.TokenPair, error) {
	if err := validators.ValidateLoginRequest(email, password); err != nil {
		uc.logger.Error("User login failed: validation error", err.Error())
		return nil, err
	}

	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		uc.logger.Error("User login failed: user not found", email)
		return nil, domainerrors.ErrInvalidCredentials
	}

	if err := uc.validateUserForLogin(user, password); err != nil {
		uc.logger.Error("User login failed: authentication failed", email)
		return nil, err
	}

	tokenPair, err := uc.authService.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		uc.logger.Error("User login failed: token generation failed", email)
		return nil, domainerrors.ErrFailedToGenerateTokens
	}

	uc.logger.Info("User logged in successfully", email)
	return tokenPair, nil
}

func (uc *authUseCase) validateUserForLogin(user *entities.User, password string) error {
	if !user.IsActive {
		return domainerrors.ErrUserDeactivated
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return domainerrors.ErrInvalidCredentials
	}

	return nil
}

func (uc *authUseCase) RefreshToken(ctx context.Context, refreshToken string) (*auth.TokenPair, error) {
	claims, err := uc.authService.ValidateToken(refreshToken)
	if err != nil {
		return nil, domainerrors.ErrInvalidToken
	}

	systemUserID := uuid.MustParse(constants.SystemUserID)
	user, err := uc.userRepo.GetByID(ctx, claims.UserID, systemUserID)
	if err != nil {
		return nil, domainerrors.ErrUserNotFound
	}

	if !user.IsActive {
		return nil, domainerrors.ErrUserAccountIsDeactivated
	}

	tokenPair, err := uc.authService.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, domainerrors.ErrFailedToGenerateTokens
	}

	return tokenPair, nil
}

func (uc *authUseCase) ValidateToken(ctx context.Context, token string) (*auth.Claims, error) {
	claims, err := uc.authService.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	if err := uc.validateUserForToken(ctx, claims.UserID); err != nil {
		return nil, err
	}

	return claims, nil
}

func (uc *authUseCase) validateUserForToken(ctx context.Context, userID uuid.UUID) error {
	systemUserID := uuid.MustParse(constants.SystemUserID)
	user, err := uc.userRepo.GetByID(ctx, userID, systemUserID)
	if err != nil {
		return domainerrors.ErrUserNotFound
	}

	if !user.IsActive {
		return domainerrors.ErrUserAccountIsDeactivated
	}

	return nil
}
