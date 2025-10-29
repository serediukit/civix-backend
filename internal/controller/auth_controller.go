package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/serediukit/civix-backend/internal/middleware"
	"github.com/serediukit/civix-backend/internal/model"
	"github.com/serediukit/civix-backend/internal/service"
	"github.com/serediukit/civix-backend/internal/util"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param input body model.CreateUserRequest true "User registration data"
// @Success 201 {object} model.UserResponse
// @Failure 400 {object} util.Response
// @Failure 409 {object} util.Response
// @Failure 500 {object} util.Response
// @Router /auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var req model.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		util.BadRequest(ctx, "Invalid request body", err)
		return
	}

	user, err := c.authService.Register(ctx.Request.Context(), &req)
	if err != nil {
		util.BadRequest(ctx, "Failed to create user", err)
		return
	}

	util.Created(ctx, user)
}

// Login handles user login
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body model.LoginRequest true "Login credentials"
// @Success 200 {object} model.Token
// @Failure 400 {object} util.Response
// @Failure 401 {object} util.Response
// @Failure 500 {object} util.Response
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req model.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		util.BadRequest(ctx, "Invalid request body", err)
		return
	}

	token, err := c.authService.Login(ctx.Request.Context(), req.Email, req.Password)
	if err != nil {
		util.Unauthorized(ctx, "Invalid credentials", err)
		return
	}

	util.Success(ctx, token)
}

// Logout handles user logout
// @Summary User logout
// @Description Invalidate the current JWT token
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} util.Response
// @Failure 401 {object} util.Response
// @Failure 500 {object} util.Response
// @Router /auth/logout [post]
func (c *AuthController) Logout(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		util.Unauthorized(ctx, "Authorization header is required", nil)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if err := c.authService.Logout(ctx.Request.Context(), tokenString); err != nil {
		util.InternalServerError(ctx, "Failed to logout", err)
		return
	}

	util.Success(ctx, gin.H{"message": "Successfully logged out"})
}

// RefreshToken handles token refresh
// @Summary Refresh JWT token
// @Description Get a new access token using a refresh token
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} model.Token
// @Failure 400 {object} util.Response
// @Failure 401 {object} util.Response
// @Failure 500 {object} util.Response
// @Router /auth/refresh [post]
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		util.Unauthorized(ctx, "Authorization header is required", nil)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := c.authService.RefreshToken(ctx.Request.Context(), tokenString)
	if err != nil {
		util.Unauthorized(ctx, "Failed to refresh token", err)
		return
	}

	util.Success(ctx, token)
}

// GetMe returns the current user's profile
// @Summary Get current user profile
// @Description Get the profile of the currently authenticated user
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} model.UserResponse
// @Failure 401 {object} util.Response
// @Failure 404 {object} util.Response
// @Failure 500 {object} util.Response
// @Router /auth/me [get]
func (c *AuthController) GetMe(ctx *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(ctx.Request.Context())
	if !exists {
		util.Unauthorized(ctx, "User not authenticated", nil)
		return
	}

	user, err := c.authService.GetProfile(ctx.Request.Context(), userID)
	if err != nil {
		util.NotFound(ctx, "User not found", err)
		return
	}

	util.Success(ctx, user)
}
