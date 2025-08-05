// Package errors provides domain-specific error types and error handling utilities
// HTTP status codes, and consistent error messaging.
package errors

import (
	"fmt"
	"net/http"
)

type ErrorCategory string

const (
	CategoryValidation   ErrorCategory = "validation"
	CategoryNotFound     ErrorCategory = "not_found"
	CategoryUnauthorized ErrorCategory = "unauthorized"
	CategoryForbidden    ErrorCategory = "forbidden"
	CategoryConflict     ErrorCategory = "conflict"
	CategoryInternal     ErrorCategory = "internal"
	CategoryDatabase     ErrorCategory = "database"
)

type AppError struct {
	Category ErrorCategory `json:"category"`
	Code     string        `json:"code"`
	Message  string        `json:"message"`
	Status   int           `json:"status"`
	Cause    error         `json:"-"`
}

func (e AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s:%s] %s: %v", e.Category, e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s:%s] %s", e.Category, e.Code, e.Message)
}

func (e AppError) Unwrap() error {
	return e.Cause
}

func NewValidationError(code, message string) *AppError {
	return &AppError{
		Category: CategoryValidation,
		Code:     code,
		Message:  message,
		Status:   http.StatusBadRequest,
	}
}

func NewNotFoundError(code, message string) *AppError {
	return &AppError{
		Category: CategoryNotFound,
		Code:     code,
		Message:  message,
		Status:   http.StatusNotFound,
	}
}

func NewUnauthorizedError(code, message string) *AppError {
	return &AppError{
		Category: CategoryUnauthorized,
		Code:     code,
		Message:  message,
		Status:   http.StatusUnauthorized,
	}
}

func NewForbiddenError(code, message string) *AppError {
	return &AppError{
		Category: CategoryForbidden,
		Code:     code,
		Message:  message,
		Status:   http.StatusForbidden,
	}
}

func NewConflictError(code, message string) *AppError {
	return &AppError{
		Category: CategoryConflict,
		Code:     code,
		Message:  message,
		Status:   http.StatusConflict,
	}
}

func NewInternalError(code, message string, cause error) *AppError {
	return &AppError{
		Category: CategoryInternal,
		Code:     code,
		Message:  message,
		Cause:    cause,
		Status:   http.StatusInternalServerError,
	}
}

func NewDatabaseError(code, message string, cause error) *AppError {
	return &AppError{
		Category: CategoryDatabase,
		Code:     code,
		Message:  message,
		Cause:    cause,
		Status:   http.StatusInternalServerError,
	}
}

var (
	ErrInvalidRequest      = NewValidationError("INVALID_REQUEST", "invalid request")
	ErrInvalidCredentials  = NewValidationError("INVALID_CREDENTIALS", "invalid credentials")
	ErrInvalidEmail        = NewValidationError("INVALID_EMAIL", "invalid email format")
	ErrInvalidID           = NewValidationError("INVALID_ID", "invalid ID")
	ErrInvalidUserID       = NewValidationError("INVALID_USER_ID", "invalid user ID")
	ErrInvalidProductID    = NewValidationError("INVALID_PRODUCT_ID", "invalid product ID")
	ErrEmailIsRequired     = NewValidationError("EMAIL_REQUIRED", "email is required")
	ErrFirstNameIsRequired = NewValidationError("FIRST_NAME_REQUIRED", "first name is required")
	ErrLastNameIsRequired  = NewValidationError("LAST_NAME_REQUIRED", "last name is required")
	ErrRoleIsRequired      = NewValidationError("ROLE_REQUIRED", "role is required")
	ErrInvalidRole         = NewValidationError("INVALID_ROLE", "invalid role")
	ErrCategoryRequired    = NewValidationError("CATEGORY_REQUIRED", "category is required")
	ErrPasswordRequired    = NewValidationError("PASSWORD_REQUIRED", "password is required")
	ErrPasswordTooShort    = NewValidationError("PASSWORD_TOO_SHORT", "password must be at least 6 characters")

	// Not found errors
	ErrUserNotFound    = NewNotFoundError("USER_NOT_FOUND", "user not found")
	ErrProductNotFound = NewNotFoundError("PRODUCT_NOT_FOUND", "product not found")

	// Unauthorized errors
	ErrInvalidOrExpiredToken       = NewUnauthorizedError("INVALID_TOKEN", "invalid or expired token")
	ErrAuthorizationHeaderRequired = NewUnauthorizedError("AUTH_HEADER_REQUIRED", "authorization header required")
	ErrUserIDNotFound              = NewUnauthorizedError("USER_ID_NOT_FOUND", "user ID not found")
	ErrUserRoleNotFound            = NewUnauthorizedError("USER_ROLE_NOT_FOUND", "user role not found")
	ErrFailedToValidateToken       = NewUnauthorizedError("TOKEN_VALIDATION_FAILED", "failed to validate token")
	ErrFailedToParseToken          = NewUnauthorizedError("TOKEN_PARSE_FAILED", "failed to parse token")
	ErrInvalidToken                = NewUnauthorizedError("INVALID_TOKEN", "invalid token")
	ErrUnexpectedSigningMethod     = NewUnauthorizedError("UNEXPECTED_SIGNING_METHOD", "unexpected signing method")
	ErrUserAccountIsDeactivated    = NewUnauthorizedError("USER_DEACTIVATED", "user account is deactivated")

	// Forbidden errors
	ErrInsufficientPermissions = NewForbiddenError("INSUFFICIENT_PERMISSIONS", "insufficient permissions")

	// Conflict errors
	ErrUserAlreadyExists    = NewConflictError("USER_EXISTS", "user already exists")
	ErrProductAlreadyExists = NewConflictError("PRODUCT_EXISTS", "product already exists")

	// Internal errors
	ErrFailedToCreateUser           = NewInternalError("USER_CREATE_FAILED", "failed to create user", nil)
	ErrFailedToUpdateUser           = NewInternalError("USER_UPDATE_FAILED", "failed to update user", nil)
	ErrFailedToDeleteUser           = NewInternalError("USER_DELETE_FAILED", "failed to delete user", nil)
	ErrFailedToGetUser              = NewInternalError("USER_GET_FAILED", "failed to get user", nil)
	ErrFailedToListUsers            = NewInternalError("USER_LIST_FAILED", "failed to list users", nil)
	ErrFailedToCreateProduct        = NewInternalError("PRODUCT_CREATE_FAILED", "failed to create product", nil)
	ErrFailedToUpdateProduct        = NewInternalError("PRODUCT_UPDATE_FAILED", "failed to update product", nil)
	ErrFailedToDeleteProduct        = NewInternalError("PRODUCT_DELETE_FAILED", "failed to delete product", nil)
	ErrFailedToGetProduct           = NewInternalError("PRODUCT_GET_FAILED", "failed to get product", nil)
	ErrFailedToListProducts         = NewInternalError("PRODUCT_LIST_FAILED", "failed to list products", nil)
	ErrFailedToGenerateAccessToken  = NewInternalError("ACCESS_TOKEN_FAILED", "failed to generate access token", nil)
	ErrFailedToGenerateRefreshToken = NewInternalError("REFRESH_TOKEN_FAILED", "failed to generate refresh token", nil)
	ErrFailedToProcessPassword      = NewInternalError("PASSWORD_PROCESS_FAILED", "failed to process password", nil)
	ErrFailedToGenerateTokens       = NewInternalError("TOKEN_GENERATION_FAILED", "failed to generate tokens", nil)

	// Deprecated aliases - kept for backward compatibility
	ErrDeleteUser      = ErrFailedToDeleteUser
	ErrUserDeactivated = ErrUserAccountIsDeactivated
)

type PermissionError struct {
	UserRole string
	Resource string
	Action   string
	Reason   string
	UserID   string
}

func (e PermissionError) Error() string {
	if e.UserID != "" {
		return fmt.Sprintf("permission denied: user %s with role %s cannot %s on %s - %s",
			e.UserID, e.UserRole, e.Action, e.Resource, e.Reason)
	}
	return fmt.Sprintf("permission denied: role %s cannot %s on %s - %s",
		e.UserRole, e.Action, e.Resource, e.Reason)
}

// NewPermissionError creates a new permission error with the given parameters
func NewPermissionError(userRole, resource, action, reason string) *PermissionError {
	return &PermissionError{
		UserRole: userRole,
		Resource: resource,
		Action:   action,
		Reason:   reason,
	}
}

// NewPermissionErrorWithUserID creates a new permission error with user ID included
func NewPermissionErrorWithUserID(userID, userRole, resource, action, reason string) *PermissionError {
	return &PermissionError{
		UserID:   userID,
		UserRole: userRole,
		Resource: resource,
		Action:   action,
		Reason:   reason,
	}
}

type RoleNotFoundError struct {
	Role string
}

func (e RoleNotFoundError) Error() string {
	return fmt.Sprintf("role not found: %s", e.Role)
}

// NewRoleNotFoundError creates a new role not found error
func NewRoleNotFoundError(role string) *RoleNotFoundError {
	return &RoleNotFoundError{Role: role}
}

type InvalidPermissionError struct {
	Permission string
	Details    string
}

func (e InvalidPermissionError) Error() string {
	return fmt.Sprintf("invalid permission: %s - %s", e.Permission, e.Details)
}

// NewInvalidPermissionError creates a new invalid permission error
func NewInvalidPermissionError(permission, details string) *InvalidPermissionError {
	return &InvalidPermissionError{
		Permission: permission,
		Details:    details,
	}
}
