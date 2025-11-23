package controller

//
// import (
// 	"github.com/gin-gonic/gin"
// 	"github.com/serediukit/civix-backend/internal/middleware"
// 	"github.com/serediukit/civix-backend/internal/model"
// 	"github.com/serediukit/civix-backend/internal/services/user"
// 	"github.com/serediukit/civix-backend/internal/util"
// )
//
// type UserController struct {
// 	userService *user.UserService
// }
//
// func NewUserController(userService *user.UserService) *UserController {
// 	return &UserController{
// 		userService: userService,
// 	}
// }
//
// // GetProfile returns the current user's profile
// // @Summary Get user profile
// // @Description Get the profile of the currently authenticated user
// // @Tags users
// // @Security BearerAuth
// // @Produce json
// // @Success 200 {object} model.UserResponse
// // @Failure 401 {object} util.Response
// // @Failure 404 {object} util.Response
// // @Failure 500 {object} util.Response
// // @Router /users/me [get]
// func (c *UserController) GetProfile(ctx *gin.Context) {
// 	userID, exists := middleware.GetUserIDFromContext(ctx.Request.Context())
// 	if !exists {
// 		util.Unauthorized(ctx, "User not authenticated", nil)
// 		return
// 	}
//
// 	user, err := c.userService.GetProfile(ctx.Request.Context(), userID)
// 	if err != nil {
// 		util.NotFound(ctx, "User not found", err)
// 		return
// 	}
//
// 	util.Success(ctx, user)
// }
//
// // UpdateProfile updates the current user's profile
// // @Summary Update user profile
// // @Description Update the profile of the currently authenticated user
// // @Tags users
// // @Security BearerAuth
// // @Accept json
// // @Produce json
// // @Param input body model.UpdateUserRequest true "User update data"
// // @Success 200 {object} model.UserResponse
// // @Failure 400 {object} util.Response
// // @Failure 401 {object} util.Response
// // @Failure 404 {object} util.Response
// // @Failure 409 {object} util.Response
// // @Failure 500 {object} util.Response
// // @Router /users/me [put]
// func (c *UserController) UpdateProfile(ctx *gin.Context) {
// 	userID, exists := middleware.GetUserIDFromContext(ctx.Request.Context())
// 	if !exists {
// 		util.Unauthorized(ctx, "User not authenticated", nil)
// 		return
// 	}
//
// 	var req model.UpdateUserRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		util.BadRequest(ctx, "Invalid request body", err)
// 		return
// 	}
//
// 	user, err := c.userService.UpdateProfile(ctx.Request.Context(), userID, &req)
// 	if err != nil {
// 		util.BadRequest(ctx, "Failed to update profile", err)
// 		return
// 	}
//
// 	util.Success(ctx, user)
// }
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
