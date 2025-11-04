package user

import (
	"context"
	"errors"

	"github.com/serediukit/civix-backend/internal/model"
	"github.com/serediukit/civix-backend/internal/repository"
)

type UserService interface {
	GetProfile(ctx context.Context, userID uint) (*model.UserResponse, error)
	UpdateProfile(ctx context.Context, userID uint, req *model.UpdateUserRequest) (*model.UserResponse, error)
	ChangePassword(ctx context.Context, userID uint, currentPassword, newPassword string) error
	DeleteAccount(ctx context.Context, userID uint) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetProfile(ctx context.Context, userID uint) (*model.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &model.UserResponse{
		ID:        user.UserID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID uint, req *model.UpdateUserRequest) (*model.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Update fields if they are provided
	if req.Name != "" {
		user.Name = req.Name
	}

	if req.Email != "" && req.Email != user.Email {
		// Check if new email is already taken
		existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
		if err == nil && existingUser != nil && existingUser.UserID != userID {
			return nil, errors.New("email already in use")
		}
		user.Email = req.Email
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, errors.New("failed to update user")
	}

	return &model.UserResponse{
		ID:        user.UserID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *userService) ChangePassword(ctx context.Context, userID uint, currentPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify current password
	if err := user.CheckPassword(currentPassword); err != nil {
		return errors.New("current password is incorrect")
	}

	// Update password (hashing is handled by BeforeSave hook)
	user.Password = newPassword

	if err := s.userRepo.Update(ctx, user); err != nil {
		return errors.New("failed to update password")
	}

	return nil
}

func (s *userService) DeleteAccount(ctx context.Context, userID uint) error {
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return errors.New("failed to delete account")
	}
	return nil
}
