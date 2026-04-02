package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"user-service/internal/config"
	"user-service/internal/handler"
	"user-service/internal/middleware"
	"user-service/internal/repository"
	"user-service/internal/service"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger := setupLogger(cfg.LogLevel)
	
	// Initialize database
	db, err := repository.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := repository.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize Redis cache
	cache := repository.NewRedisCache(cfg.RedisURL)
	defer cache.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db, cache)

	// Initialize services
	userService := service.NewUserService(userRepo, cfg.JWTSecret, cfg.JWTExpiry)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)

	// Setup router
	gin.SetMode(cfg.GinMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger(logger))
	router.Use(middleware.CORS())
	router.Use(middleware.RequestID())
	router.Use(middleware.Metrics())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes
		v1.POST("/users/register", userHandler.Register)
		v1.POST("/users/login", userHandler.Login)
		v1.POST("/users/refresh-token", userHandler.RefreshToken)
		v1.POST("/users/verify-email", userHandler.VerifyEmail)

		// Protected routes
		auth := v1.Group("/")
		auth.Use(middleware.JWTAuth(cfg.JWTSecret))
		{
			auth.GET("/users/profile", userHandler.GetProfile)
			auth.PUT("/users/profile", userHandler.UpdateProfile)
			auth.POST("/users/change-password", userHandler.ChangePassword)
			auth.DELETE("/users/account", userHandler.DeleteAccount)
		}
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Printf("Starting user service on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Println("Server exited")
}

func setupLogger(level string) *log.Logger {
	return log.New(os.Stdout, "[USER-SERVICE] ", log.LstdFlags|log.Lshortfile)
}
