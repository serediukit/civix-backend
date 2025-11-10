package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/serediukit/civix-backend/internal/contracts"
	"github.com/serediukit/civix-backend/internal/services/auth"
	"github.com/serediukit/civix-backend/pkg/hash"
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
	var req contracts.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		hash.BadRequest(ctx, "Invalid request body", err)
		return
	}

	user, err := c.authService.Register(ctx.Request.Context(), &req)
	if err != nil {
		hash.BadRequest(ctx, "Failed to create user", err)
		return
	}

	hash.Created(ctx, user)
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req contracts.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		hash.BadRequest(ctx, "Invalid request body", err)
		return
	}

	resp, err := c.authService.Login(ctx.Request.Context(), &req)
	if err != nil {
		hash.Unauthorized(ctx, "Invalid credentials", err)
		return
	}

	hash.Success(ctx, resp)
}

//
//	func (c *AuthController) Logout(ctx *gin.Context) {
//		authHeader := ctx.GetHeader("Authorization")
//		if authHeader == "" {
//			util.Unauthorized(ctx, "Authorization header is required", nil)
//			return
//		}
//
//		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
//		if err := c.authService.Logout(ctx.Request.Context(), tokenString); err != nil {
//			util.InternalServerError(ctx, "Failed to logout", err)
//			return
//		}
//
//		util.Success(ctx, gin.H{"message": "Successfully logged out"})
//	}
//func (c *AuthController) RefreshToken(ctx *gin.Context) {
//	authHeader := ctx.GetHeader("Authorization")
//	if authHeader == "" {
//		util.Unauthorized(ctx, "Authorization header is required", nil)
//		return
//	}
//
//	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
//	token, err := c.authService.RefreshToken(ctx.Request.Context(), tokenString)
//	if err != nil {
//		util.Unauthorized(ctx, "Failed to refresh token", err)
//		return
//	}
//
//	util.Success(ctx, token)
//}
