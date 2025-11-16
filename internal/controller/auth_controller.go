package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/serediukit/civix-backend/internal/contracts"
	"github.com/serediukit/civix-backend/internal/services/auth"
	"github.com/serediukit/civix-backend/pkg/util/response"
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
		response.BadRequest(ctx, "Invalid request body", err)
		return
	}

	user, err := c.authService.Register(ctx.Request.Context(), &req)
	if err != nil {
		response.BadRequest(ctx, "Failed to create user", err)
		return
	}

	response.Created(ctx, user)
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req contracts.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request body", err)
		return
	}

	resp, err := c.authService.Login(ctx.Request.Context(), &req)
	if err != nil {
		response.Unauthorized(ctx, "Invalid credentials", err)
		return
	}

	response.Success(ctx, resp)
}

func (c *AuthController) Logout(ctx *gin.Context) {
	var req contracts.LogoutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request body", err)
		return
	}

	if err := c.authService.Logout(ctx.Request.Context(), &req); err != nil {
		response.InternalServerError(ctx, "Failed to logout", err)
		return
	}

	response.Success(ctx, gin.H{"message": "Successfully logged out"})
}

func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var req contracts.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request body", err)
		return
	}

	resp, err := c.authService.RefreshToken(ctx.Request.Context(), &req)
	if err != nil {
		response.Unauthorized(ctx, "Failed to refresh token", err)
		return
	}

	response.Success(ctx, resp)
}
