package services

import (
	"context"
	"errors"
	"time"

	"github.com/serediukit/civix-backend/internal/db"

	"github.com/serediukit/civix-backend/internal/contracts"
	"github.com/serediukit/civix-backend/internal/model"
	"github.com/serediukit/civix-backend/internal/repository"
	"github.com/serediukit/civix-backend/pkg/hash"
	"github.com/serediukit/civix-backend/pkg/jwt"
	"github.com/serediukit/civix-backend/pkg/util/timeutil"
)

type AuthService interface {
	Register(ctx context.Context, req *contracts.RegisterRequest) (*contracts.RegisterResponse, error)
	Login(ctx context.Context, req *contracts.LoginRequest) (*contracts.LoginResponse, error)
	Logout(ctx context.Context, req *contracts.LogoutRequest) error
	RefreshToken(ctx context.Context, req *contracts.RefreshTokenRequest) (*contracts.RefreshTokenResponse, error)
}

type authService struct {
	userRepo   repository.UserRepository
	cityRepo   repository.CityRepository
	cachedRepo repository.CacheRepository
	jwt        *jwt.JWT
}

func NewAuthService(
	userRepo repository.UserRepository,
	cityRepo repository.CityRepository,
	cachedRepo repository.CacheRepository,
	jwt *jwt.JWT,
) AuthService {
	return &authService{
		userRepo:   userRepo,
		cityRepo:   cityRepo,
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

	city, err := s.cityRepo.GetCityByLocation(ctx, req.Location)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			city = &model.City{CityID: "661cc9c4-9cb2-48c8-9833-2aa21fd37798"} // Kyiv city_id
		} else {
			return nil, err
		}
	}

	user := &model.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Name:         req.Name,
		Surname:      req.Surname,
		PhoneNumber:  req.PhoneNumber,
		RegCityID:    city.CityID,
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return &contracts.RegisterResponse{
		Email:   user.Email,
		Name:    user.Name,
		RegTime: timeutil.Now(),
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

func (s *authService) Logout(ctx context.Context, req *contracts.LogoutRequest) error {
	errAccess := s.removeToken(ctx, req.AccessToken)
	errRefresh := s.removeToken(ctx, req.RefreshToken)

	if errAccess != nil || errRefresh != nil {
		return errors.New("failed to logout")
	}
	return nil
}

func (s *authService) RefreshToken(ctx context.Context, req *contracts.RefreshTokenRequest) (*contracts.RefreshTokenResponse, error) {
	claims, err := s.jwt.ValidateToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	blacklisted, err := s.cachedRepo.IsBlacklisted(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	if blacklisted {
		return nil, errors.New("token has been revoked")
	}

	accessTokenString, err := s.jwt.GenerateAccessToken(claims.UserID, claims.Email)
	if err != nil {
		return nil, err
	}

	accessToken, _ := s.jwt.ValidateToken(accessTokenString)

	refreshTokenString, err := s.jwt.GenerateRefreshToken(claims.UserID, claims.Email)
	if err != nil {
		return nil, err
	}

	refreshToken, _ := s.jwt.ValidateToken(refreshTokenString)

	if err = s.cachedRepo.SetBlacklist(ctx, req.RefreshToken, timeutil.Month); err != nil {
		return nil, err
	}

	return &contracts.RefreshTokenResponse{
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

func (s *authService) removeToken(ctx context.Context, token string) error {
	claims, err := s.jwt.ValidateToken(token)
	if err != nil {
		return nil
	}

	expiresIn := time.Until(claims.ExpiresAt.Time)

	return s.cachedRepo.SetBlacklist(ctx, token, expiresIn)
}
