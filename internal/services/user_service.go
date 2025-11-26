package services

import (
	"context"
	"fmt"

	"github.com/serediukit/civix-backend/internal/contracts"
	"github.com/serediukit/civix-backend/internal/middleware"
	"github.com/serediukit/civix-backend/internal/model"
	"github.com/serediukit/civix-backend/internal/repository"
)

type UserService interface {
	GetUser(ctx context.Context, req *contracts.GetUserRequest) (*contracts.GetUserResponse, error)
	UpdateProfile(ctx context.Context, req *contracts.UpdateUserRequest) (*contracts.UpdateUserResponse, error)
	// ChangePassword(ctx context.Context, userID uint, currentPassword, newPassword string) error
	// DeleteAccount(ctx context.Context, userID uint) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetUser(ctx context.Context, req *contracts.GetUserRequest) (*contracts.GetUserResponse, error) {
	var (
		user *model.User
		err  error
	)

	if req.UserID != 0 {
		user, err = s.userRepo.GetUserByID(ctx, req.UserID)
		if err != nil {
			return nil, err
		}
	}

	if req.Email != "" {
		user, err = s.userRepo.GetUserByEmail(ctx, req.Email)
		if err != nil {
			return nil, err
		}
	}

	if user == nil {
		err = fmt.Errorf("user not found for request: %+v", req)
	}

	return &contracts.GetUserResponse{
		User: user,
	}, err
}

func (s *userService) UpdateProfile(ctx context.Context, req *contracts.UpdateUserRequest) (*contracts.UpdateUserResponse, error) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("user not found in token: %+v", req)
	}

	user := &model.User{
		UserID:    userID,
		Name:      req.Name,
		Surname:   req.Surname,
		AvatarUrl: req.AvatarURL,
	}

	err := s.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &contracts.UpdateUserResponse{
		User: user,
	}, nil
}

// func (s *userService) ChangePassword(ctx context.Context, userID uint, currentPassword, newPassword string) error {
// 	user, err := s.userRepo.FindByID(ctx, userID)
// 	if err != nil {
// 		return errors.New("user not found")
// 	}
//
// 	// Verify current password
// 	if err := user.CheckPassword(currentPassword); err != nil {
// 		return errors.New("current password is incorrect")
// 	}
//
// 	// Update password (hashing is handled by BeforeSave hook)
// 	user.Password = newPassword
//
// 	if err := s.userRepo.Update(ctx, user); err != nil {
// 		return errors.New("failed to update password")
// 	}
//
// 	return nil
// }
//
// func (s *userService) DeleteAccount(ctx context.Context, userID uint) error {
// 	if err := s.userRepo.Delete(ctx, userID); err != nil {
// 		return errors.New("failed to delete account")
// 	}
// 	return nil
// }
