package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"clean-architecture-api/internal/domain/constants"
	domainerrors "clean-architecture-api/internal/domain/errors"
	"clean-architecture-api/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BaseHandler struct {
	logger logger.Logger
}

func NewBaseHandler(logger logger.Logger) *BaseHandler {
	return &BaseHandler{logger: logger}
}

func (h *BaseHandler) ParseUUID(c *gin.Context, paramName string) (uuid.UUID, error) {
	idStr := c.Param(paramName)
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, domainerrors.ErrInvalidID
	}
	return id, nil
}

func (h *BaseHandler) ParsePagination(c *gin.Context) (limit, offset int) {
	limitStr := c.DefaultQuery("limit", strconv.Itoa(constants.DefaultLimit))
	offsetStr := c.DefaultQuery("offset", strconv.Itoa(constants.DefaultOffset))

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = constants.DefaultLimit
	}

	offset, err = strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = constants.DefaultOffset
	}

	return limit, offset
}

func (h *BaseHandler) SendErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	h.logger.Error(message, err)

	// Check if it's our structured error
	var appErr *domainerrors.AppError
	if errors.As(err, &appErr) {
		c.JSON(h.getStatusCodeFromCategory(appErr.Category), gin.H{
			"error": gin.H{
				"category": appErr.Category,
				"code":     appErr.Code,
				"message":  appErr.Message,
			},
		})
		return
	}

	// Fallback for non-structured errors
	c.JSON(statusCode, gin.H{"error": err.Error()})
}

func (h *BaseHandler) getStatusCodeFromCategory(category domainerrors.ErrorCategory) int {
	switch category {
	case domainerrors.CategoryValidation:
		return http.StatusBadRequest
	case domainerrors.CategoryNotFound:
		return http.StatusNotFound
	case domainerrors.CategoryUnauthorized:
		return http.StatusUnauthorized
	case domainerrors.CategoryForbidden:
		return http.StatusForbidden
	case domainerrors.CategoryConflict:
		return http.StatusConflict
	case domainerrors.CategoryInternal, domainerrors.CategoryDatabase:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func (h *BaseHandler) SendSuccessResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{
		"success": true,
		"data":    data,
	})
}

func (h *BaseHandler) SendBadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{"error": message})
}

func (h *BaseHandler) SendNotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, gin.H{"error": message})
}

func (h *BaseHandler) SendInternalServerError(c *gin.Context, message string, err error) {
	h.logger.Error(message, err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}
