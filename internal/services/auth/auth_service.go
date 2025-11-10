package auth

import (
	"context"
	"errors"
	"github.com/serediukit/civix-backend/internal/contracts"
	"github.com/serediukit/civix-backend/internal/model"
	"github.com/serediukit/civix-backend/internal/repository"
	"github.com/serediukit/civix-backend/pkg/hash"
	"github.com/serediukit/civix-backend/pkg/jwt"
)

type AuthService interface {
	Register(ctx context.Context, req *contracts.RegisterRequest) (*contracts.RegisterResponse, error)
	Login(ctx context.Context, req *contracts.LoginRequest) (*contracts.LoginResponse, error)
	// Logout(ctx context.Context, tokenString string) error
	// RefreshToken(ctx context.Context, tokenString string) (*model.Token, error)
}

type authService struct {
	userRepo   repository.UserRepository
	cachedRepo repository.CacheRepository
	jwt        *jwt.JWT
}

func NewAuthService(
	userRepo repository.UserRepository,
	cachedRepo repository.CacheRepository,
	jwt *jwt.JWT,
) AuthService {
	return &authService{
		userRepo:   userRepo,
		cachedRepo: cachedRepo,
		jwt:        jwt,
	}
}

func (s *authService) Register(ctx context.Context, req *contracts.RegisterRequest) (*contracts.RegisterResponse, error) {
	existingUser, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := hash.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Name:         req.Name,
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return &contracts.RegisterResponse{
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *authService) Login(ctx context.Context, req *contracts.LoginRequest) (*contracts.LoginResponse, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("user with this email does not exist")
	}

	if err = hash.CheckHash(req.Password, user.PasswordHash); err != nil {
		return nil, errors.New("incorrect email or password")
	}

	accessTokenString, err := s.jwt.GenerateAccessToken(user.UserID, user.Email)
	if err != nil {
		return nil, err
	}

	accessToken, _ := s.jwt.ValidateToken(accessTokenString)

	refreshTokenString, err := s.jwt.GenerateRefreshToken(user.UserID, user.Email)
	if err != nil {
		return nil, err
	}

	refreshToken, _ := s.jwt.ValidateToken(refreshTokenString)

	return &contracts.LoginResponse{
		AccessToken: model.Token{
			Token:     accessTokenString,
			ExpiresAt: accessToken.ExpiresAt.Unix(),
		},
		RefreshToken: model.Token{
			Token:     refreshTokenString,
			ExpiresAt: refreshToken.ExpiresAt.Unix(),
		},
	}, nil
}

//
// func (s *authService) Logout(ctx context.Context, tokenString string) error {
// 	// Add token to blacklist
// 	claims, err := s.jwt.ValidateToken(tokenString)
// 	if err != nil {
// 		return err
// 	}
//
// 	expiresIn := time.Until(claims.ExpiresAt.Time)
// 	return s.cachedRepo.SetBlacklist(ctx, tokenString, expiresIn)
// }
//
// func (s *authService) RefreshToken(ctx context.Context, tokenString string) (*model.Token, error) {
// 	// Validate the refresh token
// 	claims, err := s.jwt.ValidateToken(tokenString)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	// Check if token is blacklisted
// 	blacklisted, err := s.cachedRepo.IsBlacklisted(ctx, tokenString)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if blacklisted {
// 		return nil, errors.New("token has been revoked")
// 	}
//
// 	// Generate new access token
// 	newTokenString, err := s.jwt.GenerateToken(claims.UserID, claims.Email)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	// Get new token expiration time
// 	newToken, _ := s.jwt.ValidateToken(newTokenString)
//
// 	// Add old token to blacklist
// 	s.cachedRepo.SetBlacklist(ctx, tokenString)
//
// 	return &model.Token{
// 		AccessToken: newTokenString,
// 		ExpiresAt:   newToken.ExpiresAt.Unix(),
// 	}, nil
// }
