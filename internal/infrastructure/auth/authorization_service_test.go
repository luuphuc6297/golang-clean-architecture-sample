package auth

import (
	"context"
	"testing"

	"clean-architecture-api/internal/domain/constants"
	"clean-architecture-api/internal/domain/entities"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPolicyEngine struct {
	mock.Mock
}

func (m *MockPolicyEngine) Evaluate(ctx context.Context, req *entities.PermissionRequest) (*entities.PermissionResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*entities.PermissionResponse), args.Error(1)
}

func (m *MockPolicyEngine) LoadPolicies(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockPolicyEngine) AddPolicy(ctx context.Context, policy *entities.PolicyDocument) error {
	args := m.Called(ctx, policy)
	return args.Error(0)
}

func (m *MockPolicyEngine) RemovePolicy(ctx context.Context, policyID uuid.UUID) error {
	args := m.Called(ctx, policyID)
	return args.Error(0)
}

func (m *MockPolicyEngine) GetPoliciesForRole(ctx context.Context, role string) ([]*entities.PolicyDocument, error) {
	args := m.Called(ctx, role)
	return args.Get(0).([]*entities.PolicyDocument), args.Error(1)
}

func TestNewAuthorizationService(t *testing.T) {
	mockEngine := &MockPolicyEngine{}
	service := NewAuthorizationService(mockEngine)
	assert.NotNil(t, service)
}

func TestAuthorizationService_CheckPermission(t *testing.T) {
	mockEngine := &MockPolicyEngine{}
	service := NewAuthorizationService(mockEngine)
	userID := uuid.New()

	tests := []struct {
		name          string
		userRole      string
		resource      string
		action        string
		mockResponse  *entities.PermissionResponse
		mockError     error
		expectedError bool
	}{
		{
			name:     "admin can create user",
			userRole: constants.RoleAdmin,
			resource: constants.PermissionUserCreate,
			action:   constants.ActionCreate,
			mockResponse: &entities.PermissionResponse{
				Allowed: true,
				Reason:  "allowed by policy",
			},
			mockError:     nil,
			expectedError: false,
		},
		{
			name:     "user cannot create user",
			userRole: constants.RoleUser,
			resource: constants.PermissionUserCreate,
			action:   constants.ActionCreate,
			mockResponse: &entities.PermissionResponse{
				Allowed: false,
				Reason:  "denied by policy",
			},
			mockError:     nil,
			expectedError: true,
		},
		{
			name:     "missing role should fail",
			userRole: "",
			resource: constants.PermissionUserCreate,
			action:   constants.ActionCreate,
			mockResponse: &entities.PermissionResponse{
				Allowed: false,
				Reason:  "no role",
			},
			mockError:     nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx context.Context
			if tt.userRole != "" {
				ctx = context.WithValue(context.Background(), constants.ContextUserRole, tt.userRole)
			} else {
				ctx = context.Background()
			}

			if tt.userRole != "" {
				mockEngine.On("Evaluate", mock.Anything, mock.AnythingOfType("*entities.PermissionRequest")).
					Return(tt.mockResponse, tt.mockError).Once()
			}

			err := service.CheckPermission(ctx, userID, tt.resource, tt.action)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.userRole != "" {
				mockEngine.AssertExpectations(t)
			}
		})
	}
}

func TestAuthorizationService_CheckResourcePermission(t *testing.T) {
	mockEngine := &MockPolicyEngine{}
	service := NewAuthorizationService(mockEngine)
	userID := uuid.New()
	resourceID := "test-resource-id"

	ctx := context.WithValue(context.Background(), constants.ContextUserRole, constants.RoleUser)

	mockEngine.On("Evaluate", mock.Anything, mock.MatchedBy(func(req *entities.PermissionRequest) bool {
		return req.ResourceID == resourceID
	})).Return(&entities.PermissionResponse{
		Allowed: true,
		Reason:  "allowed by policy",
	}, nil).Once()

	err := service.CheckResourcePermission(ctx, userID, constants.PermissionUserRead, constants.ActionRead, resourceID)
	assert.NoError(t, err)
	mockEngine.AssertExpectations(t)
}

func TestAuthorizationService_GetUserPermissions(t *testing.T) {
	mockEngine := &MockPolicyEngine{}
	service := NewAuthorizationService(mockEngine)
	userID := uuid.New()

	mockPolicies := []*entities.PolicyDocument{
		{
			ID:   uuid.New(),
			Name: "test-policy",
			Statements: []entities.PolicyStatement{
				{
					Effect:   constants.PolicyEffectAllow,
					Resource: constants.PermissionUserRead,
					Action:   constants.ActionRead,
				},
			},
		},
	}

	ctx := context.WithValue(context.Background(), constants.ContextUserRole, constants.RoleUser)

	mockEngine.On("GetPoliciesForRole", ctx, constants.RoleUser).
		Return(mockPolicies, nil).Once()

	permissions, err := service.GetUserPermissions(ctx, userID)
	assert.NoError(t, err)
	assert.Len(t, permissions, 1)
	assert.Equal(t, constants.PermissionUserRead, permissions[0].Resource)
	assert.Equal(t, constants.ActionRead, permissions[0].Action)

	mockEngine.AssertExpectations(t)
}

func TestAuthorizationService_ValidateRole(t *testing.T) {
	mockEngine := &MockPolicyEngine{}
	service := NewAuthorizationService(mockEngine)

	tests := []struct {
		name    string
		role    string
		wantErr bool
	}{
		{
			name:    "valid admin role",
			role:    constants.RoleAdmin,
			wantErr: false,
		},
		{
			name:    "valid user role",
			role:    constants.RoleUser,
			wantErr: false,
		},
		{
			name:    "invalid role",
			role:    "invalid_role",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateRole(tt.role)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthorizationService_QuickCheck(t *testing.T) {
	mockEngine := &MockPolicyEngine{}
	service := NewAuthorizationService(mockEngine)

	mockEngine.On("Evaluate", mock.Anything, mock.AnythingOfType("*entities.PermissionRequest")).
		Return(&entities.PermissionResponse{
			Allowed: true,
			Reason:  "allowed by policy",
		}, nil).Once()

	result := service.QuickCheck(constants.RoleAdmin, constants.PermissionUserCreate, constants.ActionCreate)
	assert.True(t, result)
	mockEngine.AssertExpectations(t)
}
