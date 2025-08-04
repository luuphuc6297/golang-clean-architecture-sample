package http

import (
	"fmt"

	"clean-architecture-api/internal/delivery/http/handlers"
	"clean-architecture-api/internal/delivery/middleware"
	"clean-architecture-api/internal/infrastructure/auth"
	"clean-architecture-api/internal/infrastructure/repository"
	"clean-architecture-api/internal/usecase"
	"clean-architecture-api/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	router *gin.Engine
	db     *gorm.DB
	logger logger.Logger
}

func NewServer(db *gorm.DB, logger logger.Logger) (*Server, error) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	server := &Server{
		router: router,
		db:     db,
		logger: logger,
	}

	if err := server.setupRoutes(); err != nil {
		return nil, fmt.Errorf("failed to setup routes: %w", err)
	}
	return server, nil
}

func (s *Server) setupRoutes() error {
	// Initialize dependencies
	handlers, authMiddleware, err := s.initializeDependencies()
	if err != nil {
		return err
	}

	// Setup health check
	s.setupHealthCheck()

	// Setup API routes
	s.setupAPIRoutes(handlers, authMiddleware)

	return nil
}

// initializeDependencies initializes all services, repositories, use cases and handlers
func (s *Server) initializeDependencies() (*routeHandlers, *middleware.AuthMiddleware, error) {
	authService, err := auth.NewAuthService()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create auth service: %w", err)
	}
	authLogger := auth.NewAuditLogger(s.logger)

	policyRepo := repository.NewPolicySQLiteRepository(s.db, s.logger)
	policyEngine := auth.NewPolicyEngine(policyRepo, s.logger)
	authzService := auth.NewAuthorizationService(policyEngine)

	userRepo := repository.NewUserRepository(s.db, authzService, authLogger, s.logger)
	productRepo := repository.NewProductRepository(s.db, authzService, authLogger, s.logger)

	authUseCase := usecase.NewAuthUseCase(userRepo, authService, s.logger)
	userUseCase := usecase.NewUserUseCase(userRepo, s.logger)
	productUseCase := usecase.NewProductUseCase(productRepo, s.logger)

	handlers := &routeHandlers{
		auth:    handlers.NewAuthHandler(authUseCase, s.logger),
		user:    handlers.NewUserHandler(userUseCase, s.logger),
		product: handlers.NewProductHandler(productUseCase, s.logger),
	}

	authMiddleware := middleware.NewAuthMiddleware(authUseCase, authzService, s.logger)

	return handlers, authMiddleware, nil
}

// routeHandlers holds all route handlers
type routeHandlers struct {
	auth    *handlers.AuthHandler
	user    *handlers.UserHandler
	product *handlers.ProductHandler
}

// setupHealthCheck sets up health check endpoint
func (s *Server) setupHealthCheck() {
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}

// setupAPIRoutes sets up all API routes
func (s *Server) setupAPIRoutes(h *routeHandlers, authMiddleware *middleware.AuthMiddleware) {
	api := s.router.Group("/api/v1")
	{
		s.setupAuthRoutes(api, h.auth)
		s.setupUserRoutes(api, h.user, authMiddleware)
		s.setupProductRoutes(api, h.product, authMiddleware)
	}
}

// setupAuthRoutes sets up authentication routes
func (s *Server) setupAuthRoutes(api *gin.RouterGroup, authHandler *handlers.AuthHandler) {
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
	}
}

// setupUserRoutes sets up user management routes
func (s *Server) setupUserRoutes(api *gin.RouterGroup, userHandler *handlers.UserHandler, authMiddleware *middleware.AuthMiddleware) {
	users := api.Group("/users")
	{
		usersProtected := users.Group("")
		usersProtected.Use(authMiddleware.UserListAccess())
		{
			usersProtected.GET("", userHandler.ListUsers)
		}

		usersProtected.Use(authMiddleware.UserReadAccess())
		{
			usersProtected.GET("/:id", userHandler.GetUserByID)
		}

		usersProtected.Use(authMiddleware.UserUpdateAccess())
		{
			usersProtected.PUT("/:id", userHandler.UpdateUser)
		}

		usersProtected.Use(authMiddleware.UserDeleteAccess())
		{
			usersProtected.DELETE("/:id", userHandler.DeleteUser)
		}
	}
}

// setupProductRoutes sets up product management routes
func (s *Server) setupProductRoutes(api *gin.RouterGroup, productHandler *handlers.ProductHandler, authMiddleware *middleware.AuthMiddleware) {
	products := api.Group("/products")
	{
		products.GET("", productHandler.ListProducts)
		products.GET("/:id", productHandler.GetProductByID)
		products.GET("/category/:category", productHandler.GetProductsByCategory)

		productsProtected := products.Group("")
		productsProtected.Use(authMiddleware.ProductCreateAccess())
		{
			productsProtected.POST("", productHandler.CreateProduct)
		}

		productsProtected.Use(authMiddleware.ProductUpdateAccess())
		{
			productsProtected.PUT("/:id", productHandler.UpdateProduct)
		}

		productsProtected.Use(authMiddleware.ProductDeleteAccess())
		{
			productsProtected.DELETE("/:id", productHandler.DeleteProduct)
		}
	}
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
