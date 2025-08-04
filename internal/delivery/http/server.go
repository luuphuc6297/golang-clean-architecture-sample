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
	authService, err := auth.NewAuthService()
	if err != nil {
		return fmt.Errorf("failed to create auth service: %w", err)
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

	authHandler := handlers.NewAuthHandler(authUseCase, s.logger)
	userHandler := handlers.NewUserHandler(userUseCase, s.logger)
	productHandler := handlers.NewProductHandler(productUseCase, s.logger)

	authMiddleware := middleware.NewAuthMiddleware(authUseCase, authzService, s.logger)

	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := s.router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

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
	return nil
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
