// Package middleware provides HTTP middleware for authentication and authorization.
// It includes middleware for token validation, role-based access control, and resource protection.
package middleware

import (
	"clean-architecture-api/internal/domain/constants"
	"clean-architecture-api/internal/domain/errors"
	"clean-architecture-api/internal/domain/repositories"
	"clean-architecture-api/internal/usecase"
	"clean-architecture-api/pkg/logger"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthMiddleware provides authentication and authorization middleware
type AuthMiddleware struct {
	authUseCase usecase.AuthUseCase
	authService repositories.AuthorizationService
	logger      logger.Logger
}

// NewAuthMiddleware creates a new authentication middleware instance
func NewAuthMiddleware(
	authUseCase usecase.AuthUseCase,
	authService repositories.AuthorizationService,
	logger logger.Logger,
) *AuthMiddleware {
	return &AuthMiddleware{
		authUseCase: authUseCase,
		authService: authService,
		logger:      logger,
	}
}

// AuthRequired middleware ensures the request has a valid authentication token
func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": errors.ErrAuthorizationHeaderRequired.Error()})
			c.Abort()
			return
		}

		claims, err := m.authUseCase.ValidateToken(c.Request.Context(), token)
		if err != nil {
			m.logger.Error(errors.ErrFailedToValidateToken.Error(), err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": errors.ErrInvalidOrExpiredToken.Error()})
			c.Abort()
			return
		}

		c.Set(string(constants.ContextUserID), claims.UserID)
		c.Set(string(constants.ContextUserEmail), claims.Email)
		c.Set(string(constants.ContextUserRole), claims.Role)

		enrichedCtx := m.authService.CreateEnrichedContext(
			c.Request.Context(),
			claims.UserID,
			claims.Role,
			claims.Email,
		)
		enrichedCtx = context.WithValue(enrichedCtx, constants.ContextClientIP, c.ClientIP())
		c.Request = c.Request.WithContext(enrichedCtx)

		c.Next()
	}
}

func (m *AuthMiddleware) ResourceAccess(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		m.AuthRequired()(c)
		if c.IsAborted() {
			return
		}

		userID, exists := c.Get(string(constants.ContextUserID))
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrUserIDNotFound.Error()})
			c.Abort()
			return
		}

		userUUID := userID.(uuid.UUID)

		if err := m.authService.CheckPermission(c.Request.Context(), userUUID, resource, action); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrInsufficientPermissions.Error()})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *AuthMiddleware) ResourceAccessWithID(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		m.AuthRequired()(c)
		if c.IsAborted() {
			return
		}

		userID, exists := c.Get(string(constants.ContextUserID))
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrUserIDNotFound.Error()})
			c.Abort()
			return
		}

		userUUID := userID.(uuid.UUID)
		resourceID := c.Param("id")

		if err := m.authService.CheckResourcePermission(c.Request.Context(), userUUID, resource, action, resourceID); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrInsufficientPermissions.Error()})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *AuthMiddleware) RoleRequired(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		m.AuthRequired()(c)
		if c.IsAborted() {
			return
		}

		userRole, exists := c.Get(string(constants.ContextUserRole))
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrUserRoleNotFound.Error()})
			c.Abort()
			return
		}

		if userRole != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": errors.ErrInsufficientPermissions.Error()})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *AuthMiddleware) AdminRequired() gin.HandlerFunc {
	return m.RoleRequired(constants.RoleAdmin)
}

func (m *AuthMiddleware) UserCreateAccess() gin.HandlerFunc {
	return m.ResourceAccess(constants.PermissionUserCreate, constants.ActionCreate)
}

func (m *AuthMiddleware) UserReadAccess() gin.HandlerFunc {
	return m.ResourceAccessWithID(constants.PermissionUserRead, constants.ActionRead)
}

func (m *AuthMiddleware) UserUpdateAccess() gin.HandlerFunc {
	return m.ResourceAccessWithID(constants.PermissionUserUpdate, constants.ActionUpdate)
}

func (m *AuthMiddleware) UserDeleteAccess() gin.HandlerFunc {
	return m.ResourceAccessWithID(constants.PermissionUserDelete, constants.ActionDelete)
}

func (m *AuthMiddleware) UserListAccess() gin.HandlerFunc {
	return m.ResourceAccess(constants.PermissionUserList, constants.ActionList)
}

func (m *AuthMiddleware) ProductCreateAccess() gin.HandlerFunc {
	return m.ResourceAccess(constants.PermissionProductCreate, constants.ActionCreate)
}

func (m *AuthMiddleware) ProductReadAccess() gin.HandlerFunc {
	return m.ResourceAccessWithID(constants.PermissionProductRead, constants.ActionRead)
}

func (m *AuthMiddleware) ProductUpdateAccess() gin.HandlerFunc {
	return m.ResourceAccessWithID(constants.PermissionProductUpdate, constants.ActionUpdate)
}

func (m *AuthMiddleware) ProductDeleteAccess() gin.HandlerFunc {
	return m.ResourceAccessWithID(constants.PermissionProductDelete, constants.ActionDelete)
}

func (m *AuthMiddleware) ProductListAccess() gin.HandlerFunc {
	return m.ResourceAccess(constants.PermissionProductList, constants.ActionList)
}

func extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return ""
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	return token
}
