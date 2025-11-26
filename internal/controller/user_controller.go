package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/serediukit/civix-backend/internal/contracts"
	"github.com/serediukit/civix-backend/internal/services"
	"github.com/serediukit/civix-backend/pkg/util/response"
)

type UserController interface {
	GetUser(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
}

type userController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return &userController{
		userService: userService,
	}
}

func (c *userController) GetUser(ctx *gin.Context) {
	var req contracts.GetUserRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, "Invalid query parameters", err)
		return
	}

	resp, err := c.userService.GetUser(ctx.Request.Context(), &req)
	if err != nil {
		response.InternalServerError(ctx, "Failed to get user", err)
		return
	}

	response.Success(ctx, resp)
}

func (c *userController) UpdateUser(ctx *gin.Context) {
	var req contracts.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "Invalid request body", err)
		return
	}

	resp, err := c.userService.UpdateProfile(ctx.Request.Context(), &req)
	if err != nil {
		response.BadRequest(ctx, "Failed to update profile", err)
		return
	}

	response.Success(ctx, resp)
}

//
// // ChangePassword changes the current user's password
// // @Summary Change password
// // @Description Change the password of the currently authenticated user
// // @Tags users
// // @Security BearerAuth
// // @Accept json
// // @Produce json
// // @Param input body changePasswordRequest true "Password change data"
// // @Success 200 {object} util.Response
// // @Failure 400 {object} util.Response
// // @Failure 401 {object} util.Response
// // @Failure 500 {object} util.Response
// // @Router /users/me/password [put]
// func (c *UserController) ChangePassword(ctx *gin.Context) {
// 	userID, exists := middleware.GetUserIDFromContext(ctx.Request.Context())
// 	if !exists {
// 		util.Unauthorized(ctx, "User not authenticated", nil)
// 		return
// 	}
//
// 	var req changePasswordRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		util.BadRequest(ctx, "Invalid request body", err)
// 		return
// 	}
//
// 	if err := c.userService.ChangePassword(ctx.Request.Context(), userID, req.CurrentPassword, req.NewPassword); err != nil {
// 		util.BadRequest(ctx, "Failed to change password", err)
// 		return
// 	}
//
// 	util.Success(ctx, gin.H{"message": "Password updated successfully"})
// }
//
// // DeleteAccount deletes the current user's account
// // @Summary Delete account
// // @Description Delete the currently authenticated user's account
// // @Tags users
// // @Security BearerAuth
// // @Produce json
// // @Success 200 {object} util.Response
// // @Failure 401 {object} util.Response
// // @Failure 500 {object} util.Response
// // @Router /users/me [delete]
// func (c *UserController) DeleteAccount(ctx *gin.Context) {
// 	userID, exists := middleware.GetUserIDFromContext(ctx.Request.Context())
// 	if !exists {
// 		util.Unauthorized(ctx, "User not authenticated", nil)
// 		return
// 	}
//
// 	if err := c.userService.DeleteAccount(ctx.Request.Context(), userID); err != nil {
// 		util.InternalServerError(ctx, "Failed to delete account", err)
// 		return
// 	}
//
// 	util.Success(ctx, gin.H{"message": "Account deleted successfully"})
// }
//
// type changePasswordRequest struct {
// 	CurrentPassword string `json:"current_password" binding:"required"`
// 	NewPassword     string `json:"new_password" binding:"required,min=6"`
// }
