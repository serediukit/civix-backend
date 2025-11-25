package server

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/serediukit/civix-backend/internal/config"
	"github.com/serediukit/civix-backend/internal/controller"
	"github.com/serediukit/civix-backend/internal/middleware"
	"github.com/serediukit/civix-backend/internal/repository"
	"github.com/serediukit/civix-backend/internal/services"
	"github.com/serediukit/civix-backend/pkg/database"
	"github.com/serediukit/civix-backend/pkg/jwt"
	"github.com/serediukit/civix-backend/pkg/redis"
	"github.com/sirupsen/logrus"
)

type Server struct {
	config *config.Config
	router *gin.Engine
	logger *logrus.Logger
}

func NewServer(config *config.Config) *Server {
	return &Server{config: config}
}

func (s *Server) Run() error {
	s.logger = logrus.New()

	// Configure logger
	var formatter logrus.Formatter
	if s.config.Server.GinMode == gin.ReleaseMode {
		// Production: compact JSON
		formatter = &logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		}
		s.logger.SetLevel(logrus.InfoLevel)
	} else {
		// Development: pretty-printed JSON with indentation
		formatter = &logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
			PrettyPrint:     true,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		}
		s.logger.SetLevel(logrus.DebugLevel)
	}
	s.logger.SetFormatter(formatter)

	s.logger.SetOutput(os.Stdout)

	// Set Gin mode
	gin.SetMode(s.config.Server.GinMode)

	ctx := context.Background()

	// Initialize database connections
	store, err := database.NewDB(ctx, s.config.GetDBConfig())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer store.Close()

	// Initialize Redis
	redisClient, err := redis.NewRedis(s.config.GetRedisConfig())
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	userRepo := repository.NewUserRepository(store)
	reportRepo := repository.NewReportRepository(store)
	cityRepo := repository.NewCityRepository(store)
	cacheRepo := repository.NewCacheRepository(redisClient)

	// Initialize utilities
	jwtTokenizer := jwt.NewJWT(s.config.GetJWTConfig())

	// Initialize services
	authService := services.NewAuthService(userRepo, cityRepo, cacheRepo, jwtTokenizer)
	// userService := user.NewUserService(userRepo)
	reportService := services.NewReportService(reportRepo, cityRepo)

	// Initialize controllers
	authController := controller.NewAuthController(authService)
	// userController := controller.NewUserController(userService)
	reportController := controller.NewReportController(reportService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtTokenizer, cacheRepo)

	// Create router
	s.router = setupRouter(authController, reportController, authMiddleware, s.logger)

	srv := &http.Server{
		Addr:    ":" + s.config.Server.Port,
		Handler: s.router,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	go func() {
		log.Printf("Server is running on port %s\n", s.config.Server.Port)

		if err = srv.ListenAndServeTLS("cert.pem", "key.pem"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()

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
