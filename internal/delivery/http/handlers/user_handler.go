package handlers

import (
	"clean-architecture-api/internal/domain/constants"
	"clean-architecture-api/internal/domain/entities"
	domainerrors "clean-architecture-api/internal/domain/errors"
	"clean-architecture-api/internal/usecase"
	"clean-architecture-api/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	*BaseHandler
	userUseCase usecase.UserUseCase
}

func NewUserHandler(userUseCase usecase.UserUseCase, logger logger.Logger) *UserHandler {
	return &UserHandler{
		BaseHandler: NewBaseHandler(logger),
		userUseCase: userUseCase,
	}
}

type UpdateUserRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Role      string `json:"role" binding:"required"`
	IsActive  bool   `json:"is_active"`
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	targetUserID, err := h.ParseUUID(c, "id")
	if err != nil {
		h.SendErrorResponse(c, 0, "Invalid user ID", err)
		return
	}

	currentUserID := h.getCurrentUserID(c)
	user, err := h.userUseCase.GetByID(c.Request.Context(), targetUserID, currentUserID)
	if err != nil {
		h.SendErrorResponse(c, 0, "Failed to get user", err)
		return
	}

	h.SendSuccessResponse(c, http.StatusOK, gin.H{"user": user})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	targetUserID, err := h.ParseUUID(c, "id")
	if err != nil {
		h.SendErrorResponse(c, 0, "Invalid user ID", err)
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErr := domainerrors.NewValidationError("INVALID_REQUEST_BODY", "request body validation failed")
		h.SendErrorResponse(c, 0, "Invalid request", validationErr)
		return
	}

	user := h.createUserFromRequest(targetUserID, req)
	currentUserID := h.getCurrentUserID(c)

	if err := h.userUseCase.Update(c.Request.Context(), user, currentUserID); err != nil {
		h.SendErrorResponse(c, 0, "Failed to update user", err)
		return
	}

	h.SendSuccessResponse(c, http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (h *UserHandler) createUserFromRequest(userID uuid.UUID, req UpdateUserRequest) *entities.User {
	return &entities.User{
		BaseEntity: entities.BaseEntity{
			ID: userID,
		},
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      req.Role,
		IsActive:  req.IsActive,
	}
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	targetUserID, err := h.ParseUUID(c, "id")
	if err != nil {
		h.SendErrorResponse(c, 0, "Invalid user ID", err)
		return
	}

	currentUserID := h.getCurrentUserID(c)
	if err := h.userUseCase.Delete(c.Request.Context(), targetUserID, currentUserID); err != nil {
		h.SendErrorResponse(c, 0, "Failed to delete user", err)
		return
	}

	h.SendSuccessResponse(c, http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	limit, offset := h.ParsePagination(c)
	currentUserID := h.getCurrentUserID(c)

	users, err := h.userUseCase.List(c.Request.Context(), limit, offset, currentUserID)
	if err != nil {
		h.SendInternalServerError(c, "Failed to list users", err)
		return
	}

	h.SendSuccessResponse(c, http.StatusOK, gin.H{"users": users})
}

func (h *UserHandler) getCurrentUserID(c *gin.Context) uuid.UUID {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(uuid.UUID); ok {
			return id
		}
	}
	return uuid.MustParse(constants.SystemUserID)
}
