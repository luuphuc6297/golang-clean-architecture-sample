package usecase

import (
	"clean-architecture-api/internal/domain/entities"
	"clean-architecture-api/internal/infrastructure/auth"
	"clean-architecture-api/pkg/logger"
	"context"
	"testing"

	domainerrors "clean-architecture-api/internal/domain/errors"

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

type RegisterTestCase struct {
	name          string
	email         string
	password      string
	firstName     string
	lastName      string
	setupMocks    func(*MockUserRepository, *MockAuthService, *MockLogger)
	expectedError error
}

// getRegisterSuccessTestCase returns the success test case for Register function
func getRegisterSuccessTestCase() RegisterTestCase {
	return RegisterTestCase{
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
			mockLogger.On("Info", mock.Anything, mock.Anything).Return()
		},
		expectedError: nil,
	}
}

// getRegisterValidationTestCases returns validation failure test cases
func getRegisterValidationTestCases() []RegisterTestCase {
	return []RegisterTestCase{
		{
			name:      "Failure - Invalid email",
			email:     "invalid-email",
			password:  "password123",
			firstName: "John",
			lastName:  "Doe",
			setupMocks: func(mockRepo *MockUserRepository, mockAuth *MockAuthService, mockLogger *MockLogger) {
				// No GetByEmail call expected since validation fails early
				mockLogger.On("Error", mock.Anything, mock.Anything).Return()
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
				// No GetByEmail call expected since validation fails early
				mockLogger.On("Error", mock.Anything, mock.Anything).Return()
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
				// No GetByEmail call expected since validation fails early
				mockLogger.On("Error", mock.Anything, mock.Anything).Return()
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
				mockLogger.On("Error", mock.Anything, mock.Anything).Return()
			},
			expectedError: domainerrors.ErrFailedToCreateUser,
		},
	}
}

// getRegisterConflictTestCase returns user conflict test case
func getRegisterConflictTestCase() RegisterTestCase {
	existingUser := &entities.User{
		Email:     "existing@example.com",
		FirstName: "Existing",
		LastName:  "User",
	}

	return RegisterTestCase{
		name:      "Failure - User already exists",
		email:     "existing@example.com",
		password:  "password123",
		firstName: "John",
		lastName:  "Doe",
		setupMocks: func(mockRepo *MockUserRepository, mockAuth *MockAuthService, mockLogger *MockLogger) {
			// Mock GetByEmail to return existing user
			mockRepo.On("GetByEmail", mock.Anything, "existing@example.com").Return(existingUser, nil)
			mockLogger.On("Error", mock.Anything, mock.Anything).Return()
		},
		expectedError: domainerrors.ErrUserAlreadyExists,
	}
}

// getRegisterDatabaseTestCases returns database-related failure test cases
func getRegisterDatabaseTestCases() []RegisterTestCase {
	return []RegisterTestCase{
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
				mockLogger.On("Error", mock.Anything, mock.Anything).Return()
			},
			expectedError: domainerrors.ErrFailedToCreateUser,
		},
	}
}

// runRegisterTest executes a single register test case
func runRegisterTest(t *testing.T, tt RegisterTestCase) {
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
}

// getRegisterTestCases returns all test cases for Register function
func getRegisterTestCases() []RegisterTestCase {
	var tests []RegisterTestCase
	tests = append(tests, getRegisterSuccessTestCase())
	tests = append(tests, getRegisterConflictTestCase())
	tests = append(tests, getRegisterValidationTestCases()...)
	tests = append(tests, getRegisterDatabaseTestCases()...)
	return tests
}

func TestAuthUseCase_Register(t *testing.T) {
	tests := getRegisterTestCases()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runRegisterTest(t, tt)
		})
	}
}

// setupLoginTestData creates test data for login tests
func setupLoginTestData(t *testing.T) (*entities.User, *auth.TokenPair, uuid.UUID) {
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

	return validUser, validTokenPair, validUserID
}

type LoginTestCase struct {
	name          string
	email         string
	password      string
	setupMocks    func(*MockUserRepository, *MockAuthService, *MockLogger)
	expectedToken *auth.TokenPair
	expectedError error
}

// getLoginSuccessTestCase returns the success test case for Login function
func getLoginSuccessTestCase(validUser *entities.User, validTokenPair *auth.TokenPair, validUserID uuid.UUID) LoginTestCase {
	return LoginTestCase{
		name:     "Success - Valid login",
		email:    "test@example.com",
		password: "password123",
		setupMocks: func(mockRepo *MockUserRepository, mockAuth *MockAuthService, mockLogger *MockLogger) {
			// Mock GetByEmail to return valid user
			mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(validUser, nil)
			// Mock GenerateTokenPair to return valid token pair
			mockAuth.On("GenerateTokenPair", validUserID, "test@example.com", "user").Return(validTokenPair, nil)
			mockLogger.On("Info", mock.Anything, mock.Anything).Return()
		},
		expectedToken: validTokenPair,
		expectedError: nil,
	}
}

// getLoginFailureTestCases returns failure test cases for Login function
func getLoginFailureTestCases(validUser *entities.User, validTokenPair *auth.TokenPair, validUserID uuid.UUID) []LoginTestCase {
	return []LoginTestCase{
		{
			name:     "Failure - User not found",
			email:    "nonexistent@example.com",
			password: "password123",
			setupMocks: func(mockRepo *MockUserRepository, mockAuth *MockAuthService, mockLogger *MockLogger) {
				// Mock GetByEmail to return error
				mockRepo.On("GetByEmail", mock.Anything, "nonexistent@example.com").Return(nil, domainerrors.ErrUserNotFound)
				mockLogger.On("Error", mock.Anything, mock.Anything).Return()
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
				mockLogger.On("Error", mock.Anything, mock.Anything).Return()
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
				mockLogger.On("Error", mock.Anything, mock.Anything).Return()
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
				mockLogger.On("Error", mock.Anything, mock.Anything).Return()
			},
			expectedToken: nil,
			expectedError: domainerrors.ErrFailedToGenerateTokens,
		},
	}
}

// getLoginTestCases returns all test cases for Login function
func getLoginTestCases(validUser *entities.User, validTokenPair *auth.TokenPair, validUserID uuid.UUID) []LoginTestCase {
	var tests []LoginTestCase
	tests = append(tests, getLoginSuccessTestCase(validUser, validTokenPair, validUserID))
	tests = append(tests, getLoginFailureTestCases(validUser, validTokenPair, validUserID)...)
	return tests
}

// runLoginTest executes a single login test case
func runLoginTest(t *testing.T, tt LoginTestCase) {
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
}

func TestAuthUseCase_Login(t *testing.T) {
	validUser, validTokenPair, validUserID := setupLoginTestData(t)
	tests := getLoginTestCases(validUser, validTokenPair, validUserID)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runLoginTest(t, tt)
		})
	}
}
