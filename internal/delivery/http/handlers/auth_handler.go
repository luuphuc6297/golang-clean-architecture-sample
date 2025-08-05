package handlers

import (
	"clean-architecture-api/internal/domain/errors"
	"clean-architecture-api/internal/usecase"
	"clean-architecture-api/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	*BaseHandler
	authUseCase usecase.AuthUseCase
}

// NewAuthHandler creates a new authentication handler instance
func NewAuthHandler(authUseCase usecase.AuthUseCase, logger logger.Logger) *AuthHandler {
	return &AuthHandler{
		BaseHandler: NewBaseHandler(logger),
		authUseCase: authUseCase,
	}
}

type RegisterRequest struct {
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendBadRequest(c, errors.ErrInvalidRequest.Error())
		return
	}

	user, err := h.authUseCase.Register(c.Request.Context(), req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		h.SendErrorResponse(c, http.StatusBadRequest, "Registration failed", err)
		return
	}

	h.SendSuccessResponse(c, http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendBadRequest(c, errors.ErrInvalidRequest.Error())
		return
	}

	tokenPair, err := h.authUseCase.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		h.SendErrorResponse(c, http.StatusUnauthorized, "Login failed", err)
		return
	}

	h.SendSuccessResponse(c, http.StatusOK, gin.H{
		"message": "Login successful",
		"tokens":  tokenPair,
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendBadRequest(c, errors.ErrInvalidRequest.Error())
		return
	}

	tokenPair, err := h.authUseCase.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		h.SendErrorResponse(c, http.StatusUnauthorized, "Token refresh failed", err)
		return
	}

	h.SendSuccessResponse(c, http.StatusOK, gin.H{
		"message": "Token refreshed successfully",
		"tokens":  tokenPair,
	})
}
