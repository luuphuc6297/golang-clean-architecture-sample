package auth

import (
	"clean-architecture-api/internal/domain/constants"
	"clean-architecture-api/internal/domain/entities"
	"clean-architecture-api/internal/domain/errors"
	"clean-architecture-api/internal/domain/repositories"
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

type AuthorizationServiceImpl struct {
	policyEngine repositories.PolicyEngine
}

func NewAuthorizationService(policyEngine repositories.PolicyEngine) repositories.AuthorizationService {
	return &AuthorizationServiceImpl{
		policyEngine: policyEngine,
	}
}

func (s *AuthorizationServiceImpl) CheckPermission(ctx context.Context, userID uuid.UUID, resource, action string) error {
	return s.CheckResourcePermission(ctx, userID, resource, action, "")
}

func (s *AuthorizationServiceImpl) CheckResourcePermission(ctx context.Context, userID uuid.UUID, resource, action, resourceID string) error {
	userRole, err := s.validateUserRole(ctx)
	if err != nil {
		return err
	}

	req := &entities.PermissionRequest{
		UserID:     userID,
		Role:       userRole,
		Resource:   resource,
		Action:     action,
		ResourceID: resourceID,
		Context:    s.buildContextData(ctx, resourceID),
	}

	response, err := s.policyEngine.Evaluate(ctx, req)
	if err != nil {
		return errors.NewPermissionError(userRole, resource, action, "policy evaluation failed")
	}

	if !response.Allowed {
		return errors.NewPermissionError(userRole, resource, action, response.Reason)
	}

	return nil
}

func (s *AuthorizationServiceImpl) GetUserPermissions(ctx context.Context, _ uuid.UUID) ([]entities.Permission, error) {
	userRole, err := s.validateUserRole(ctx)
	if err != nil {
		return nil, err
	}

	policies, err := s.policyEngine.GetPoliciesForRole(ctx, userRole)
	if err != nil {
		return nil, err
	}

	return s.extractPermissionsFromPolicies(policies, userRole), nil
}

func (s *AuthorizationServiceImpl) GetEffectivePermissions(ctx context.Context, userID uuid.UUID) ([]entities.Permission, error) {
	permissions, err := s.GetUserPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}

	userRole, _ := s.validateUserRole(ctx)
	contextData := s.buildContextData(ctx, "")
	var effectivePermissions []entities.Permission

	for _, permission := range permissions {
		req := &entities.PermissionRequest{
			UserID:   userID,
			Role:     userRole,
			Resource: permission.Resource,
			Action:   permission.Action,
			Context:  contextData,
		}

		response, evalErr := s.policyEngine.Evaluate(ctx, req)
		if evalErr == nil && response.Allowed {
			effectivePermissions = append(effectivePermissions, permission)
		}
	}

	return effectivePermissions, nil
}

func (s *AuthorizationServiceImpl) QuickCheck(userRole, resource, action string) bool {
	ctx := context.WithValue(context.Background(), constants.ContextUserRole, userRole)
	userID := uuid.New()
	err := s.CheckPermission(ctx, userID, resource, action)
	return err == nil
}

func (s *AuthorizationServiceImpl) ValidateRole(userRole string) error {
	validRoles := []string{constants.RoleAdmin, constants.RoleUser}
	for _, role := range validRoles {
		if userRole == role {
			return nil
		}
	}
	return errors.NewRoleNotFoundError(userRole)
}

func (s *AuthorizationServiceImpl) GetAllowedActionsForRole(userRole, resource string) ([]string, error) {
	ctx := context.WithValue(context.Background(), constants.ContextUserRole, userRole)
	userID := uuid.New()

	permissions, err := s.GetUserPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}

	var actions []string
	for _, permission := range permissions {
		if permission.Resource == resource {
			actions = append(actions, permission.Action)
		}
	}

	return actions, nil
}

func (s *AuthorizationServiceImpl) validateUserRole(ctx context.Context) (string, error) {
	userRole, exists := ctx.Value(constants.ContextUserRole).(string)
	if !exists || userRole == "" {
		return "", errors.ErrUserRoleNotFound
	}
	return userRole, nil
}

func (s *AuthorizationServiceImpl) buildContextData(ctx context.Context, resourceID string) map[string]interface{} {
	contextData := make(map[string]interface{})

	if userID, exists := ctx.Value(constants.ContextUserID).(uuid.UUID); exists {
		contextData[string(constants.ContextUserID)] = userID.String()
	}

	if userRole, exists := ctx.Value(constants.ContextUserRole).(string); exists {
		contextData[string(constants.ContextUserRole)] = userRole
	}

	if userEmail, exists := ctx.Value(constants.ContextUserEmail).(string); exists {
		contextData[string(constants.ContextUserEmail)] = userEmail
	}

	if clientIP, exists := ctx.Value(constants.ContextClientIP).(string); exists {
		contextData[string(constants.ContextClientIP)] = clientIP
	}

	if resourceID != "" {
		contextData["resource_id"] = resourceID
	}

	return contextData
}

func (s *AuthorizationServiceImpl) extractPermissionsFromPolicies(policies []*entities.PolicyDocument, userRole string) []entities.Permission {
	var permissions []entities.Permission
	for _, policy := range policies {
		for _, statement := range policy.Statements {
			if statement.Effect == constants.PolicyEffectAllow {
				permissions = append(permissions, entities.Permission{
					Resource: statement.Resource,
					Action:   statement.Action,
					Role:     userRole,
				})
			}
		}
	}
	return permissions
}

func (s *AuthorizationServiceImpl) CreateEnrichedContext(baseCtx context.Context, userID uuid.UUID, role, email string) context.Context {
	ctx := context.WithValue(baseCtx, constants.ContextUserID, userID)
	ctx = context.WithValue(ctx, constants.ContextUserRole, role)

	if email != "" {
		ctx = context.WithValue(ctx, constants.ContextUserEmail, email)
	}

	return ctx
}

func (s *AuthorizationServiceImpl) SerializeContextForMicroservice(ctx context.Context) (string, error) {
	contextData := s.buildContextData(ctx, "")
	data, err := json.Marshal(contextData)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *AuthorizationServiceImpl) CreateContextFromMicroserviceData(baseCtx context.Context, data string) (context.Context, error) {
	var contextData map[string]interface{}
	err := json.Unmarshal([]byte(data), &contextData)
	if err != nil {
		return nil, err
	}

	ctx := baseCtx

	if userIDStr, exists := contextData[string(constants.ContextUserID)].(string); exists {
		if userID, parseErr := uuid.Parse(userIDStr); parseErr == nil {
			ctx = context.WithValue(ctx, constants.ContextUserID, userID)
		}
	}

	if userRole, exists := contextData[string(constants.ContextUserRole)].(string); exists {
		ctx = context.WithValue(ctx, constants.ContextUserRole, userRole)
	}

	if userEmail, exists := contextData[string(constants.ContextUserEmail)].(string); exists {
		ctx = context.WithValue(ctx, constants.ContextUserEmail, userEmail)
	}

	return ctx, nil
}
