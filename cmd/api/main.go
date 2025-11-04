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
	"github.com/joho/godotenv"
	"github.com/serediukit/civix-backend/internal/config"
	"github.com/serediukit/civix-backend/internal/controller"
	"github.com/serediukit/civix-backend/internal/middleware"
	"github.com/serediukit/civix-backend/internal/service"
	"github.com/serediukit/civix-backend/internal/util"
	"github.com/serediukit/civix-backend/pkg/database"
	"github.com/serediukit/civix-backend/pkg/jwt"
	"github.com/serediukit/civix-backend/pkg/redis"
)

func main() {
	// Load environment variables from .env file
	_ = godotenv.Load()

	// Initialize configuration
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set Gin mode
	gin.SetMode(config.Server.GinMode)

	ctx := context.Background()

	// Initialize database connections
	db, err := database.NewDB(ctx, config.GetDBConfig())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Redis
	redisClient, err := redis.NewRedis(config.GetRedisConfig())
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Initialize utilities
	jwt := jwt.NewJWT(config.GetJWTConfig())

	// Initialize services
	authService := service.NewAuthService(userRepo, redisRepo, config, jwt)
	userService := service.NewUserService(userRepo)

	// Initialize controllers
	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtUtil, redisRepo)

	// Create router
	router := setupRouter(authController, userController, authMiddleware)

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + config.Server.Port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Server is running on port %s\n", config.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}

func setupRouter(
	authController *controller.AuthController,
	userController *controller.UserController,
	authMiddleware *middleware.AuthMiddleware,
) *gin.Engine {
	r := gin.New()

	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// API v1
	v1 := r.Group("/api/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
			auth.POST("/logout", authMiddleware.AuthRequired(), authController.Logout)
			auth.POST("/refresh", authController.RefreshToken)
			auth.GET("/me", authMiddleware.AuthRequired(), authController.GetMe)
		}

		// User routes
		users := v1.Group("/users")
		users.Use(authMiddleware.AuthRequired())
		{
			users.GET("/me", userController.GetProfile)
			users.PUT("/me", userController.UpdateProfile)
			users.PUT("/me/password", userController.ChangePassword)
			users.DELETE("/me", userController.DeleteAccount)
		}

		// Add more routes here...
	}

	return r
}
