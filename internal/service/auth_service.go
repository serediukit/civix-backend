package service

import (
	"context"
	"errors"
	"time"

	"github.com/serediukit/civix-backend/internal/config"
	"github.com/serediukit/civix-backend/internal/model"
	"github.com/serediukit/civix-backend/internal/repository"
	"github.com/serediukit/civix-backend/internal/util"
)

type AuthService interface {
	Register(ctx context.Context, req *model.CreateUserRequest) (*model.UserResponse, error)
	Login(ctx context.Context, email, password string) (*model.Token, error)
	Logout(ctx context.Context, tokenString string) error
	RefreshToken(ctx context.Context, tokenString string) (*model.Token, error)
}

type authService struct {
	userRepo      repository.UserRepository
	redisRepo     repository.RedisRepository
	jwtUtil       *util.JWTUtil
	tokenExpiry   time.Duration
	refreshExpiry time.Duration
}

func NewAuthService(
	userRepo repository.UserRepository,
	redisRepo repository.RedisRepository,
	cfg *config.Config,
	jwtUtil *util.JWTUtil,
) AuthService {
	return &authService{
		userRepo:      userRepo,
		redisRepo:     redisRepo,
		jwtUtil:       jwtUtil,
		tokenExpiry:   cfg.JWT.Expiration,
		refreshExpiry: cfg.JWT.Expiration * 24 * 7, // 7 days
	}
}

func (s *authService) Register(ctx context.Context, req *model.CreateUserRequest) (*model.UserResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Create new user
	user := &model.User{
		Email:    req.Email,
		Password: req.Password, // Will be hashed by BeforeCreate hook
		Name:     req.Name,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return &model.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (*model.Token, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	if err := user.CheckPassword(password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate access token
	tokenString, err := s.jwtUtil.GenerateToken(user.ID, user.Email, s.tokenExpiry)
	if err != nil {
		return nil, err
	}

	// Get token expiration time
	token, _ := s.jwtUtil.ValidateToken(tokenString)

	return &model.Token{
		AccessToken: tokenString,
		ExpiresAt:   token.ExpiresAt.Unix(),
	}, nil
}

func (s *authService) Logout(ctx context.Context, tokenString string) error {
	// Add token to blacklist
	claims, err := s.jwtUtil.ValidateToken(tokenString)
	if err != nil {
		return err
	}

	expiresIn := time.Until(claims.ExpiresAt.Time)
	return s.redisRepo.SetBlacklist(ctx, tokenString, expiresIn)
}

func (s *authService) RefreshToken(ctx context.Context, tokenString string) (*model.Token, error) {
	// Validate the refresh token
	claims, err := s.jwtUtil.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Check if token is blacklisted
	blacklisted, err := s.redisRepo.IsBlacklisted(ctx, tokenString)
	if err != nil {
		return nil, err
	}
	if blacklisted {
		return nil, errors.New("token has been revoked")
	}

	// Generate new access token
	newTokenString, err := s.jwtUtil.GenerateToken(claims.UserID, claims.Email, s.tokenExpiry)
	if err != nil {
		return nil, err
	}

	// Get new token expiration time
	newToken, _ := s.jwtUtil.ValidateToken(newTokenString)

	// Add old token to blacklist
	s.redisRepo.SetBlacklist(ctx, tokenString, s.tokenExpiry)

	return &model.Token{
		AccessToken: newTokenString,
		ExpiresAt:   newToken.ExpiresAt.Unix(),
	}, nil
}
