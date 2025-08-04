package handlers

import (
	"clean-architecture-api/internal/domain/constants"
	"clean-architecture-api/internal/domain/entities"
	"clean-architecture-api/internal/domain/errors"
	"clean-architecture-api/internal/usecase"
	"clean-architecture-api/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ProductHandler handles HTTP requests for product operations
type ProductHandler struct {
	*BaseHandler
	productUseCase usecase.ProductUseCase
}

// NewProductHandler creates a new product handler instance
func NewProductHandler(productUseCase usecase.ProductUseCase, logger logger.Logger) *ProductHandler {
	return &ProductHandler{
		BaseHandler:    NewBaseHandler(logger),
		productUseCase: productUseCase,
	}
}

// CreateProductRequest represents the request body for creating a product
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"gte=0"`
	Category    string  `json:"category"`
}

// UpdateProductRequest represents the request body for updating a product
type UpdateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"gte=0"`
	Category    string  `json:"category"`
}

// CreateProduct handles product creation requests
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendBadRequest(c, errors.ErrInvalidRequest.Error())
		return
	}

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		h.SendInternalServerError(c, "Failed to get user ID", err)
		return
	}

	product := h.createProductFromRequest(req)

	if err := h.productUseCase.Create(c.Request.Context(), product, userID); err != nil {
		h.SendErrorResponse(c, http.StatusBadRequest, "Failed to create product", err)
		return
	}

	h.SendSuccessResponse(c, http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"product": product,
	})
}

func (h *ProductHandler) getUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get(string(constants.ContextUserID))
	if !exists {
		return uuid.Nil, errors.ErrUserIDNotFound
	}
	return userID.(uuid.UUID), nil
}

func (h *ProductHandler) createProductFromRequest(req CreateProductRequest) *entities.Product {
	return &entities.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
	}
}

// GetProductByID handles requests to retrieve a product by ID
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	productID, err := h.ParseUUID(c, "id")
	if err != nil {
		h.SendBadRequest(c, errors.ErrInvalidProductID.Error())
		return
	}

	product, err := h.productUseCase.GetByID(c.Request.Context(), productID)
	if err != nil {
		h.SendNotFound(c, err.Error())
		return
	}

	h.SendSuccessResponse(c, http.StatusOK, gin.H{"product": product})
}

// UpdateProduct handles product update requests
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	productID, err := h.ParseUUID(c, "id")
	if err != nil {
		h.SendBadRequest(c, errors.ErrInvalidProductID.Error())
		return
	}

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendBadRequest(c, errors.ErrInvalidRequest.Error())
		return
	}

	product := h.createProductFromRequestWithID(productID, req)

	if err := h.productUseCase.Update(c.Request.Context(), product); err != nil {
		h.SendErrorResponse(c, http.StatusBadRequest, "Failed to update product", err)
		return
	}

	h.SendSuccessResponse(c, http.StatusOK, gin.H{"message": "Product updated successfully"})
}

func (h *ProductHandler) createProductFromRequestWithID(productID uuid.UUID, req UpdateProductRequest) *entities.Product {
	return &entities.Product{
		BaseEntity: entities.BaseEntity{
			ID: productID,
		},
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
	}
}

// DeleteProduct handles product deletion requests
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	productID, err := h.ParseUUID(c, "id")
	if err != nil {
		h.SendBadRequest(c, errors.ErrInvalidProductID.Error())
		return
	}

	if err := h.productUseCase.Delete(c.Request.Context(), productID); err != nil {
		h.SendErrorResponse(c, http.StatusBadRequest, "Failed to delete product", err)
		return
	}

	h.SendSuccessResponse(c, http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// ListProducts handles requests to list all products
func (h *ProductHandler) ListProducts(c *gin.Context) {
	limit, offset := h.ParsePagination(c)

	products, err := h.productUseCase.List(c.Request.Context(), limit, offset)
	if err != nil {
		h.SendInternalServerError(c, "Failed to list products", err)
		return
	}

	h.SendSuccessResponse(c, http.StatusOK, gin.H{"products": products})
}

// GetProductsByCategory handles requests to get products by category
func (h *ProductHandler) GetProductsByCategory(c *gin.Context) {
	category := c.Param("category")
	if category == "" {
		h.SendBadRequest(c, errors.ErrCategoryRequired.Error())
		return
	}

	limit, offset := h.ParsePagination(c)

	products, err := h.productUseCase.GetByCategory(c.Request.Context(), category, limit, offset)
	if err != nil {
		h.SendInternalServerError(c, "Failed to get products by category", err)
		return
	}

	h.SendSuccessResponse(c, http.StatusOK, gin.H{"products": products})
}
