package usecase

import (
	"context"
	"testing"

	"clean-architecture-api/internal/domain/entities"
	domainerrors "clean-architecture-api/internal/domain/errors"
	"clean-architecture-api/internal/infrastructure/auth"
	"clean-architecture-api/pkg/logger"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entities.User, userID uuid.UUID) error {
	args := m.Called(ctx, user, userID)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entities.User, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) GetAll(ctx context.Context) ([]*entities.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entities.User, userID uuid.UUID) error {
	args := m.Called(ctx, user, userID)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	args := m.Called(ctx, id, deletedBy)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, limit, offset int, userID uuid.UUID) ([]*entities.User, error) {
	args := m.Called(ctx, limit, offset, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.User), args.Error(1)
}

func (m *MockUserRepository) ValidateAccess(ctx context.Context, userID uuid.UUID, action string) error {
	args := m.Called(ctx, userID, action)
	return args.Error(0)
}

func (m *MockUserRepository) AuditLog(ctx context.Context, userID uuid.UUID, action string, entity *entities.User) error {
	args := m.Called(ctx, userID, action, entity)
	return args.Error(0)
}

// Mock AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) GenerateTokenPair(userID uuid.UUID, email, role string) (*auth.TokenPair, error) {
	args := m.Called(userID, email, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.TokenPair), args.Error(1)
}

func (m *MockAuthService) ValidateToken(tokenString string) (*auth.Claims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.Claims), args.Error(1)
}

func (m *MockAuthService) RefreshTokenPair(refreshToken string) (*auth.TokenPair, error) {
	args := m.Called(refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.TokenPair), args.Error(1)
}

// Mock Logger
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(args ...any) {
	m.Called(args...)
}

func (m *MockLogger) Error(args ...any) {
	m.Called(args...)
}

func (m *MockLogger) Fatal(args ...any) {
	m.Called(args...)
}

func (m *MockLogger) Warn(args ...any) {
	m.Called(args...)
}

func (m *MockLogger) Debug(args ...any) {
	m.Called(args...)
}

func (m *MockLogger) WithField(key string, value any) logger.Logger {
	args := m.Called(key, value)
	if args.Get(0) == nil {
		return m
	}
	return args.Get(0).(logger.Logger)
}

func (m *MockLogger) WithError(err error) logger.Logger {
	args := m.Called(err)
	if args.Get(0) == nil {
		return m
	}
	return args.Get(0).(logger.Logger)
}

// Test setup helper
func setupAuthUseCaseTest() (*authUseCase, *MockUserRepository, *MockAuthService, *MockLogger) {
	mockUserRepo := &MockUserRepository{}
	mockAuthService := &MockAuthService{}
	mockLogger := &MockLogger{}

	authUC := &authUseCase{
		BaseUseCase: *NewBaseUseCase(mockLogger),
		userRepo:    mockUserRepo,
		authService: mockAuthService,
	}

	return authUC, mockUserRepo, mockAuthService, mockLogger
}

func TestAuthUseCase_Register(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		password      string
		firstName     string
		lastName      string
		setupMocks    func(*MockUserRepository, *MockAuthService, *MockLogger)
		expectedUser  *entities.User
		expectedError error
	}{
		{
			name:      "Success - Register new user",
			email:     "test@example.com",
			password:  "password123",
			firstName: "John",
			lastName:  "Doe",
			setupMocks: func(mockRepo *MockUserRepository, mockAuth *MockAuthService, mockLogger *MockLogger) {
				// Mock GetByEmail to return nil (user doesn't exist)
				mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, domainerrors.ErrUserNotFound)
				// Mock Create to return nil (success)
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.User"), mock.AnythingOfType("uuid.UUID")).Return(nil)
				mockLogger.On("Info", mock.Anything).Return()
			},
			expectedError: nil,
		},
		{
			name:      "Failure - User already exists",
			email:     "existing@example.com",
			password:  "password123",
			firstName: "John",
			lastName:  "Doe",
			setupMocks: func(mockRepo *MockUserRepository, mockAuth *MockAuthService, mockLogger *MockLogger) {
				// Mock GetByEmail to return existing user
				existingUser := &entities.User{
					Email:     "existing@example.com",
					FirstName: "Existing",
					LastName:  "User",
				}
				mockRepo.On("GetByEmail", mock.Anything, "existing@example.com").Return(existingUser, nil)
				mockLogger.On("Error", mock.Anything).Return()
			},
			expectedError: domainerrors.ErrUserAlreadyExists,
		},
		{
			name:      "Failure - Invalid email",
			email:     "invalid-email",
			password:  "password123",
			firstName: "John",
			lastName:  "Doe",
			setupMocks: func(mockRepo *MockUserRepository, mockAuth *MockAuthService, mockLogger *MockLogger) {
				// Mock GetByEmail to return nil (user doesn't exist)
				mockRepo.On("GetByEmail", mock.Anything, "invalid-email").Return(nil, domainerrors.ErrUserNotFound)
				mockLogger.On("Error", mock.Anything).Return()
			},
			expectedError: domainerrors.ErrInvalidEmail,
		},
		{
			name:      "Failure - Empty first name",
			email:     "test@example.com",
			password:  "password123",
			firstName: "",
			lastName:  "Doe",
			setupMocks: func(mockRepo *MockUserRepository, mockAuth *MockAuthService, mockLogger *MockLogger) {
				// Mock GetByEmail to return nil (user doesn't exist)
				mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, domainerrors.ErrUserNotFound)
				mockLogger.On("Error", mock.Anything).Return()
			},
			expectedError: domainerrors.ErrFirstNameIsRequired,
		},
		{
			name:      "Failure - Empty last name",
			email:     "test@example.com",
			password:  "password123",
			firstName: "John",
			lastName:  "",
			setupMocks: func(mockRepo *MockUserRepository, mockAuth *MockAuthService, mockLogger *MockLogger) {
				// Mock GetByEmail to return nil (user doesn't exist)
				mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, domainerrors.ErrUserNotFound)
				mockLogger.On("Error", mock.Anything).Return()
			},
			expectedError: domainerrors.ErrLastNameIsRequired,
		},
		{
			name:      "Failure - Database error during create",
			email:     "test@example.com",
			password:  "password123",
			firstName: "John",
			lastName:  "Doe",
			setupMocks: func(mockRepo *MockUserRepository, mockAuth *MockAuthService, mockLogger *MockLogger) {
				// Mock GetByEmail to return nil (user doesn't exist)
				mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, domainerrors.ErrUserNotFound)
				// Mock Create to return error
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.User"), mock.AnythingOfType("uuid.UUID")).Return(domainerrors.ErrFailedToCreateUser)
				mockLogger.On("Error", mock.Anything).Return()
			},
			expectedError: domainerrors.ErrFailedToCreateUser,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authUC, mockRepo, mockAuth, mockLogger := setupAuthUseCaseTest()
			tt.setupMocks(mockRepo, mockAuth, mockLogger)

			ctx := context.Background()
			user, err := authUC.Register(ctx, tt.email, tt.password, tt.firstName, tt.lastName)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.email, user.Email)
				assert.Equal(t, tt.firstName, user.FirstName)
				assert.Equal(t, tt.lastName, user.LastName)
				assert.Equal(t, "user", user.Role)
				assert.True(t, user.IsActive)
				// Password should be hashed
				assert.NotEqual(t, tt.password, user.Password)
			}

			mockRepo.AssertExpectations(t)
			mockAuth.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
		})
	}
}

func TestAuthUseCase_Login(t *testing.T) {
	validUserID := uuid.New()

	// Create real bcrypt hash for testing
	testHelper := NewTestHelper()
	hashedPassword, err := testHelper.HashPassword("password123")
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	validUser := &entities.User{
		BaseEntity: entities.BaseEntity{ID: validUserID},
		Email:      "test@example.com",
		Password:   hashedPassword,
		FirstName:  "John",
		LastName:   "Doe",
		Role:       "user",
		IsActive:   true,
	}

	validTokenPair := &auth.TokenPair{
		AccessToken:  "access_token_here",
		RefreshToken: "refresh_token_here",
		ExpiresIn:    1703123456,
	}

	tests := []struct {
		name          string
		email         string
		password      string
		setupMocks    func(*MockUserRepository, *MockAuthService, *MockLogger)
		expectedToken *auth.TokenPair
		expectedError error
	}{
		{
			name:     "Success - Valid login",
			email:    "test@example.com",
			password: "password123",
			setupMocks: func(mockRepo *MockUserRepository, mockAuth *MockAuthService, mockLogger *MockLogger) {
				// Mock GetByEmail to return valid user
				mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(validUser, nil)
				// Mock GenerateTokenPair to return valid token pair
				mockAuth.On("GenerateTokenPair", validUserID, "test@example.com", "user").Return(validTokenPair, nil)
				mockLogger.On("Info", mock.Anything).Return()
			},
			expectedToken: validTokenPair,
			expectedError: nil,
		},
		{
			name:     "Failure - User not found",
			email:    "nonexistent@example.com",
			password: "password123",
			setupMocks: func(mockRepo *MockUserRepository, mockAuth *MockAuthService, mockLogger *MockLogger) {
				// Mock GetByEmail to return error
				mockRepo.On("GetByEmail", mock.Anything, "nonexistent@example.com").Return(nil, domainerrors.ErrUserNotFound)
				mockLogger.On("Error", mock.Anything).Return()
			},
			expectedToken: nil,
			expectedError: domainerrors.ErrInvalidCredentials,
		},
		{
			name:     "Failure - User deactivated",
			email:    "deactivated@example.com",
			password: "password123",
			setupMocks: func(mockRepo *MockUserRepository, mockAuth *MockAuthService, mockLogger *MockLogger) {
				// Mock GetByEmail to return deactivated user
				deactivatedUser := &entities.User{
					BaseEntity: entities.BaseEntity{ID: validUserID},
					Email:      "deactivated@example.com",
					Password:   validUser.Password,
					FirstName:  "John",
					LastName:   "Doe",
					Role:       "user",
					IsActive:   false,
				}
				mockRepo.On("GetByEmail", mock.Anything, "deactivated@example.com").Return(deactivatedUser, nil)
				mockLogger.On("Error", mock.Anything).Return()
			},
			expectedToken: nil,
			expectedError: domainerrors.ErrUserDeactivated,
		},
		{
			name:     "Failure - Invalid password",
			email:    "test@example.com",
			password: "wrongpassword",
			setupMocks: func(mockRepo *MockUserRepository, mockAuth *MockAuthService, mockLogger *MockLogger) {
				// Mock GetByEmail to return valid user
				mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(validUser, nil)
				mockLogger.On("Error", mock.Anything).Return()
			},
			expectedToken: nil,
			expectedError: domainerrors.ErrInvalidCredentials,
		},
		{
			name:     "Failure - Token generation error",
			email:    "test@example.com",
			password: "password123",
			setupMocks: func(mockRepo *MockUserRepository, mockAuth *MockAuthService, mockLogger *MockLogger) {
				// Mock GetByEmail to return valid user
				mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(validUser, nil)
				// Mock GenerateTokenPair to return error
				mockAuth.On("GenerateTokenPair", validUserID, "test@example.com", "user").Return(nil, domainerrors.ErrFailedToGenerateTokens)
				mockLogger.On("Error", mock.Anything).Return()
			},
			expectedToken: nil,
			expectedError: domainerrors.ErrFailedToGenerateTokens,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authUC, mockRepo, mockAuth, mockLogger := setupAuthUseCaseTest()
			tt.setupMocks(mockRepo, mockAuth, mockLogger)

			ctx := context.Background()
			tokenPair, err := authUC.Login(ctx, tt.email, tt.password)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, tokenPair)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tokenPair)
				assert.Equal(t, tt.expectedToken, tokenPair)
			}

			mockRepo.AssertExpectations(t)
			mockAuth.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
		})
	}
}
