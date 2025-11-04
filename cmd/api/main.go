package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/serediukit/civix-backend/internal/repository"
	"github.com/serediukit/civix-backend/internal/services/auth"
	"github.com/serediukit/civix-backend/internal/services/user"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/serediukit/civix-backend/internal/config"
	"github.com/serediukit/civix-backend/internal/controller"
	"github.com/serediukit/civix-backend/internal/middleware"
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
	store, err := database.NewDB(ctx, config.GetDBConfig())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer store.Close()

	// Initialize Redis
	redisClient, err := redis.NewRedis(config.GetRedisConfig())
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	userRepo := repository.NewUserRepository(store)
	// reportRepo := repository.NewReportRepository(db)
	cacheRepo := repository.NewCacheRepository(redisClient)

	// Initialize utilities
	jwt := jwt.NewJWT(config.GetJWTConfig())

	// Initialize services
	authService := auth.NewAuthService(userRepo, cacheRepo, jwt)
	userService := user.NewUserService(userRepo)
	// reportService := reports.NewReportService(reportRepo)

	// Initialize controllers
	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService)
	// reportController := controller.NewReportController(reportService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwt, cacheRepo)

	// Create router
	router := setupRouter(authController, userController, authMiddleware)

	if err = startServer(config, router); err != nil {
		log.Fatalf("Failed server: %v", err)
	}

	log.Println("Server exiting")
}

func startServer(config *config.Config, router *gin.Engine) error {
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

	return nil
}
