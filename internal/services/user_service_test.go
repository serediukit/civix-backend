package services

import (
	"context"
	"testing"

	"github.com/serediukit/civix-backend/internal/contracts"
	"github.com/serediukit/civix-backend/internal/db"
	"github.com/serediukit/civix-backend/internal/middleware"
	"github.com/serediukit/civix-backend/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestUserService_GetUser_ByID_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	userService := NewUserService(mockUserRepo)

	ctx := context.Background()
	userID := uint64(1)
	expectedUser := &model.User{
		UserID:  userID,
		Email:   "test@example.com",
		Name:    "John",
		Surname: "Doe",
	}

	req := &contracts.GetUserRequest{
		UserID: userID,
	}

	mockUserRepo.On("GetUserByID", ctx, userID).Return(expectedUser, nil)

	resp, err := userService.GetUser(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedUser, resp.User)
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_GetUser_ByEmail_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	userService := NewUserService(mockUserRepo)

	ctx := context.Background()
	email := "test@example.com"
	expectedUser := &model.User{
		UserID:  1,
		Email:   email,
		Name:    "John",
		Surname: "Doe",
	}

	req := &contracts.GetUserRequest{
		Email: email,
	}

	mockUserRepo.On("GetUserByEmail", ctx, email).Return(expectedUser, nil)

	resp, err := userService.GetUser(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedUser, resp.User)
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_GetUser_UserNotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	userService := NewUserService(mockUserRepo)

	ctx := context.Background()
	userID := uint64(999)

	req := &contracts.GetUserRequest{
		UserID: userID,
	}

	mockUserRepo.On("GetUserByID", ctx, userID).Return(nil, db.ErrNotFound)

	resp, err := userService.GetUser(ctx, req)

	assert.Error(t, err)
	assert.NotNil(t, resp)
	assert.Nil(t, resp.User)
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_GetUser_NoIdentifier(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	userService := NewUserService(mockUserRepo)

	ctx := context.Background()
	req := &contracts.GetUserRequest{}

	resp, err := userService.GetUser(ctx, req)

	assert.Error(t, err)
	assert.NotNil(t, resp)
	assert.Nil(t, resp.User)
	assert.Contains(t, err.Error(), "user not found")
}

func TestUserService_UpdateProfile_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	userService := NewUserService(mockUserRepo)

	userID := uint64(1)
	ctx := middleware.SetUserIDInContext(context.Background(), userID)

	req := &contracts.UpdateUserRequest{
		Name:      "UpdatedName",
		Surname:   "UpdatedSurname",
		AvatarURL: "https://example.com/avatar.jpg",
	}

	expectedUser := &model.User{
		UserID:    userID,
		Name:      req.Name,
		Surname:   req.Surname,
		AvatarUrl: req.AvatarURL,
	}

	mockUserRepo.On("UpdateUser", ctx, expectedUser).Return(nil)

	resp, err := userService.UpdateProfile(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.Name, resp.User.Name)
	assert.Equal(t, req.Surname, resp.User.Surname)
	assert.Equal(t, req.AvatarURL, resp.User.AvatarUrl)
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_UpdateProfile_NoUserInContext(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	userService := NewUserService(mockUserRepo)

	ctx := context.Background()
	req := &contracts.UpdateUserRequest{
		Name:    "UpdatedName",
		Surname: "UpdatedSurname",
	}

	resp, err := userService.UpdateProfile(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "user not found in token")
}
