package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nshmdayo/github-copilot-sample/backend/internal/config"
	"github.com/nshmdayo/github-copilot-sample/backend/internal/handler"
	"github.com/nshmdayo/github-copilot-sample/backend/internal/middleware"
	"github.com/nshmdayo/github-copilot-sample/backend/internal/model"
	"github.com/nshmdayo/github-copilot-sample/backend/internal/repository"
	"github.com/nshmdayo/github-copilot-sample/backend/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize database
	db, err := initDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Auto migrate database schema
	if err := db.AutoMigrate(&model.User{}, &model.Todo{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	todoRepo := repository.NewTodoRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg)
	todoService := service.NewTodoService(todoRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	todoHandler := handler.NewTodoHandler(todoService)

	// Initialize Gin router
	router := setupRouter(cfg, authService, authHandler, todoHandler)

	// Start server
	address := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting server on %s", address)
	if err := router.Run(address); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func initDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.Database.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

func setupRouter(cfg *config.Config, authService service.AuthService, authHandler *handler.AuthHandler, todoHandler *handler.TodoHandler) *gin.Engine {
	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	router := gin.New()

	// Add middlewares
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.CORSMiddleware(
		cfg.CORS.AllowOrigins,
		cfg.CORS.AllowMethods,
		cfg.CORS.AllowHeaders,
		cfg.CORS.AllowCredentials,
	))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Todo API is running",
		})
	})

	// API routes
	api := router.Group("/api/v1")

	// Auth routes
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/me", middleware.AuthMiddleware(authService), authHandler.Me)
	}

	// Todo routes (protected)
	todos := api.Group("/todos")
	todos.Use(middleware.AuthMiddleware(authService))
	{
		todos.POST("", todoHandler.Create)
		todos.GET("", todoHandler.GetList)
		todos.GET("/:id", todoHandler.GetByID)
		todos.PUT("/:id", todoHandler.Update)
		todos.DELETE("/:id", todoHandler.Delete)
		todos.PATCH("/:id/toggle", todoHandler.ToggleStatus)
	}

	return router
}
