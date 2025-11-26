package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/serediukit/civix-backend/internal/controller"
	"github.com/serediukit/civix-backend/internal/middleware"
	"github.com/sirupsen/logrus"
)

func setupRouter(
	authController controller.AuthController,
	// userController *controller.UserController,
	reportController controller.ReportController,
	authMiddleware *middleware.AuthMiddleware,
	logger *logrus.Logger,
) *gin.Engine {
	r := gin.New()

	// Middleware
	r.Use(middleware.SimpleRequestLogger(logger)) // Custom detailed request/response logger
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
		}

		// // User routes
		// users := v1.Group("/users")
		// users.Use(authMiddleware.AuthRequired())
		// {
		// 	users.GET("/me", userController.GetProfile)
		// 	users.PUT("/me", userController.UpdateProfile)
		// 	users.PUT("/me/password", userController.ChangePassword)
		// 	users.DELETE("/me", userController.DeleteAccount)
		// }

		// Report routes
		reports := v1.Group("/reports")
		reports.Use(authMiddleware.AuthRequired())
		{
			reports.GET("/", reportController.GetReports)
			reports.POST("/", reportController.CreateReport)
			// reports.PUT("/:id", reportController.UpdateReport)
			// reports.DELETE("/:id", reportController.DeleteReport)
		}
	}

	return r
}
