package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/serediukit/civix-backend/internal/contracts"
	"github.com/serediukit/civix-backend/internal/services"
	"github.com/serediukit/civix-backend/pkg/util/response"
)

type AuthController interface {
	Register(router *gin.Context)
	Login(router *gin.Context)
	Logout(router *gin.Context)
	RefreshToken(router *gin.Context)
}

type authController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) AuthController {
	return &authController{
		authService: authService,
	}
}

func (c *authController) Register(ctx *gin.Context) {
	var req contracts.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request body", err)
		return
	}

	resp, err := c.authService.Register(ctx.Request.Context(), &req)
	if err != nil {
		response.BadRequest(ctx, "Failed to create user", err)
		return
	}

	response.Created(ctx, resp)
}

func (c *authController) Login(ctx *gin.Context) {
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

func (c *authController) Logout(ctx *gin.Context) {
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

func (c *authController) RefreshToken(ctx *gin.Context) {
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
