package auth

import (
	"clean-architecture-api/internal/domain/constants"
	"clean-architecture-api/internal/domain/entities"
	"clean-architecture-api/internal/domain/errors"
	"clean-architecture-api/internal/domain/repositories"
	"clean-architecture-api/pkg/logger"
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type PolicyEngineImpl struct {
	policyRepo repositories.PolicyRepository
	logger     logger.Logger
	cache      map[string][]*entities.PolicyDocument
	mutex      sync.RWMutex
}

func NewPolicyEngine(policyRepo repositories.PolicyRepository, logger logger.Logger) repositories.PolicyEngine {
	engine := &PolicyEngineImpl{
		policyRepo: policyRepo,
		logger:     logger,
		cache:      make(map[string][]*entities.PolicyDocument),
	}

	if err := engine.LoadPolicies(context.Background()); err != nil {
		logger.Error("Failed to load initial policies", err)
	}

	return engine
}

func (pe *PolicyEngineImpl) Evaluate(_ context.Context, req *entities.PermissionRequest) (*entities.PermissionResponse, error) {
	if req == nil {
		return &entities.PermissionResponse{
			Allowed: false,
			Reason:  "invalid request",
		}, errors.ErrInvalidRequest
	}

	policies := pe.getPoliciesFromCache(req.Role)
	if len(policies) == 0 {
		pe.logger.Info(fmt.Sprintf("No policies found for role: %s", req.Role))
		return &entities.PermissionResponse{
			Allowed: false,
			Reason:  "no policies found for role",
		}, nil
	}

	response := pe.evaluatePolicies(policies, req)
	pe.logEvaluation(req, response)

	return response, nil
}

func (pe *PolicyEngineImpl) evaluatePolicies(policies []*entities.PolicyDocument, req *entities.PermissionRequest) *entities.PermissionResponse {
	var allowPolicies []string
	var denyPolicies []string

	for _, policy := range policies {
		for _, statement := range policy.Statements {
			if pe.statementMatches(statement, req) {
				switch statement.Effect {
				case constants.PolicyEffectAllow:
					allowPolicies = append(allowPolicies, policy.Name)
				case constants.PolicyEffectDeny:
					denyPolicies = append(denyPolicies, policy.Name)
				}
			}
		}
	}

	if len(denyPolicies) > 0 {
		return &entities.PermissionResponse{
			Allowed:  false,
			Reason:   "denied by policy",
			Policies: denyPolicies,
		}
	}

	if len(allowPolicies) > 0 {
		return &entities.PermissionResponse{
			Allowed:  true,
			Reason:   "allowed by policy",
			Policies: allowPolicies,
		}
	}

	return &entities.PermissionResponse{
		Allowed: false,
		Reason:  "no matching policy found",
	}
}

func (pe *PolicyEngineImpl) statementMatches(statement entities.PolicyStatement, req *entities.PermissionRequest) bool {
	return pe.matchesPrincipal(statement.Principal, req.Role) &&
		pe.matchesAction(statement.Action, req.Action) &&
		pe.matchesResource(statement.Resource, req.Resource) &&
		pe.matchesConditions(statement.Conditions, req)
}

func (pe *PolicyEngineImpl) matchesPrincipal(principal, role string) bool {
	return principal == "*" || principal == "role:"+role
}

func (pe *PolicyEngineImpl) matchesAction(policyAction, requestAction string) bool {
	return policyAction == "*" || policyAction == requestAction
}

func (pe *PolicyEngineImpl) matchesResource(policyResource, requestResource string) bool {
	return policyResource == "*" || policyResource == requestResource
}

func (pe *PolicyEngineImpl) matchesConditions(conditions map[string]interface{}, req *entities.PermissionRequest) bool {
	if len(conditions) == 0 {
		return true
	}

	for key, expectedValue := range conditions {
		if key == "resource_owner" {
			if !pe.checkResourceOwnership(req) {
				return false
			}
			continue
		}

		contextValue, exists := req.Context[key]
		if !exists || contextValue != expectedValue {
			return false
		}
	}

	return true
}

// checkResourceOwnership validates resource ownership for the permission request
func (pe *PolicyEngineImpl) checkResourceOwnership(req *entities.PermissionRequest) bool {
	if req.ResourceID == "" {
		return true
	}

	contextOwner, exists := req.Context["resource_owner_id"]
	if !exists {
		return false
	}

	return contextOwner == req.UserID.String()
}

func (pe *PolicyEngineImpl) logEvaluation(req *entities.PermissionRequest, response *entities.PermissionResponse) {
	pe.logger.Info(fmt.Sprintf(
		"Policy evaluation: UserID=%s, Role=%s, Resource=%s, Action=%s, ResourceID=%s, Allowed=%t, Reason=%s",
		req.UserID.String(),
		req.Role,
		req.Resource,
		req.Action,
		req.ResourceID,
		response.Allowed,
		response.Reason,
	))
}

func (pe *PolicyEngineImpl) LoadPolicies(ctx context.Context) error {
	policies, err := pe.policyRepo.GetActive(ctx)
	if err != nil {
		return err
	}

	pe.mutex.Lock()
	defer pe.mutex.Unlock()

	pe.cache = make(map[string][]*entities.PolicyDocument)
	for _, policy := range policies {
		for _, statement := range policy.Statements {
			role := pe.extractRoleFromPrincipal(statement.Principal)
			if role != "" {
				pe.cache[role] = append(pe.cache[role], policy)
			}
		}
	}

	pe.logger.Info(fmt.Sprintf("Loaded %d policies into cache", len(policies)))
	return nil
}

func (pe *PolicyEngineImpl) extractRoleFromPrincipal(principal string) string {
	if principal == "*" {
		return "*"
	}
	if len(principal) > 5 && principal[:5] == "role:" {
		return principal[5:]
	}
	return ""
}

func (pe *PolicyEngineImpl) getPoliciesFromCache(role string) []*entities.PolicyDocument {
	pe.mutex.RLock()
	defer pe.mutex.RUnlock()

	var allPolicies []*entities.PolicyDocument

	if policies, exists := pe.cache[role]; exists {
		allPolicies = append(allPolicies, policies...)
	}

	if globalPolicies, exists := pe.cache["*"]; exists {
		allPolicies = append(allPolicies, globalPolicies...)
	}

	return pe.deduplicatePolicies(allPolicies)
}

func (pe *PolicyEngineImpl) deduplicatePolicies(policies []*entities.PolicyDocument) []*entities.PolicyDocument {
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

func (pe *PolicyEngineImpl) AddPolicy(ctx context.Context, policy *entities.PolicyDocument) error {
	if err := pe.validatePolicy(policy); err != nil {
		return err
	}

	if err := pe.policyRepo.Create(ctx, policy); err != nil {
		return err
	}

	return pe.LoadPolicies(ctx)
}

func (pe *PolicyEngineImpl) validatePolicy(policy *entities.PolicyDocument) error {
	if policy.Name == "" {
		return errors.ErrInvalidRequest
	}

	for _, statement := range policy.Statements {
		if !statement.IsValid() {
			return errors.ErrInvalidRequest
		}
	}

	return nil
}

func (pe *PolicyEngineImpl) RemovePolicy(ctx context.Context, policyID uuid.UUID) error {
	if err := pe.policyRepo.Delete(ctx, policyID); err != nil {
		return err
	}

	return pe.LoadPolicies(ctx)
}

// GetPoliciesForRole retrieves all policies for a specific role
func (pe *PolicyEngineImpl) GetPoliciesForRole(ctx context.Context, role string) ([]*entities.PolicyDocument, error) {
	return pe.policyRepo.GetByRole(ctx, role)
}
