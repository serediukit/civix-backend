package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/serediukit/civix-backend/internal/model"
	"github.com/serediukit/civix-backend/internal/services/auth"
	"github.com/serediukit/civix-backend/internal/util"
)

type AuthController struct {
	authService auth.AuthService
}

func NewAuthController(authService auth.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (c *AuthController) Register(ctx *gin.Context) {
	util.Success(ctx, "hehe")

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

// func (c *AuthController) Login(ctx *gin.Context) {
// 	var req model.LoginRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		util.BadRequest(ctx, "Invalid request body", err)
// 		return
// 	}
//
// 	token, err := c.authService.Login(ctx.Request.Context(), req.Email, req.Password)
// 	if err != nil {
// 		util.Unauthorized(ctx, "Invalid credentials", err)
// 		return
// 	}
//
// 	util.Success(ctx, token)
// }
//
// func (c *AuthController) Logout(ctx *gin.Context) {
// 	authHeader := ctx.GetHeader("Authorization")
// 	if authHeader == "" {
// 		util.Unauthorized(ctx, "Authorization header is required", nil)
// 		return
// 	}
//
// 	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
// 	if err := c.authService.Logout(ctx.Request.Context(), tokenString); err != nil {
// 		util.InternalServerError(ctx, "Failed to logout", err)
// 		return
// 	}
//
// 	util.Success(ctx, gin.H{"message": "Successfully logged out"})
// }
//
// func (c *AuthController) RefreshToken(ctx *gin.Context) {
// 	authHeader := ctx.GetHeader("Authorization")
// 	if authHeader == "" {
// 		util.Unauthorized(ctx, "Authorization header is required", nil)
// 		return
// 	}
//
// 	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
// 	token, err := c.authService.RefreshToken(ctx.Request.Context(), tokenString)
// 	if err != nil {
// 		util.Unauthorized(ctx, "Failed to refresh token", err)
// 		return
// 	}
//
// 	util.Success(ctx, token)
// }
//
// func (c *AuthController) GetMe(ctx *gin.Context) {
// 	userID, exists := middleware.GetUserIDFromContext(ctx.Request.Context())
// 	if !exists {
// 		util.Unauthorized(ctx, "User not authenticated", nil)
// 		return
// 	}
//
// 	user, err := c.authService.GetProfile(ctx.Request.Context(), userID)
// 	if err != nil {
// 		util.NotFound(ctx, "User not found", err)
// 		return
// 	}
//
// 	util.Success(ctx, user)
// }
