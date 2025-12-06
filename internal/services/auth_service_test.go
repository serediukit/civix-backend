package services

import (
	"context"
	"testing"
	"time"

	"github.com/serediukit/civix-backend/internal/contracts"
	"github.com/serediukit/civix-backend/internal/db"
	"github.com/serediukit/civix-backend/internal/model"
	"github.com/serediukit/civix-backend/pkg/hash"
	"github.com/serediukit/civix-backend/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repositories
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id uint64) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

type MockCityRepository struct {
	mock.Mock
}

func (m *MockCityRepository) GetCityByLocation(ctx context.Context, location model.Location) (*model.City, error) {
	args := m.Called(ctx, location)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.City), args.Error(1)
}

type MockCacheRepository struct {
	mock.Mock
}

func (m *MockCacheRepository) SetBlacklist(ctx context.Context, token string, duration time.Duration) error {
	args := m.Called(ctx, token, duration)
	return args.Error(0)
}

func (m *MockCacheRepository) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	args := m.Called(ctx, token)
	return args.Bool(0), args.Error(1)
}

func (m *MockCacheRepository) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func TestAuthService_Register_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockCityRepo := new(MockCityRepository)
	mockCacheRepo := new(MockCacheRepository)
	jwtService := jwt.NewJWT(&jwt.JWTConfig{
		Secret:            "test-secret",
		TokenExpiration:   15 * time.Minute,
		RefreshExpiration: 7 * 24 * time.Hour,
	})

	authService := NewAuthService(mockUserRepo, mockCityRepo, mockCacheRepo, jwtService)

	ctx := context.Background()
	req := &contracts.RegisterRequest{
		Email:       "test@example.com",
		Password:    "password123",
		Name:        "John",
		Surname:     "Doe",
		PhoneNumber: "+1234567890",
		Location: model.Location{
			Lat: 50.4501,
			Lng: 30.5234,
		},
	}

	city := &model.City{
		CityID: "123e4567-e89b-12d3-a456-426614174000",
		Name:   "Kyiv",
	}

	mockUserRepo.On("GetUserByEmail", ctx, req.Email).Return(nil, db.ErrNotFound)
	mockCityRepo.On("GetCityByLocation", ctx, req.Location).Return(city, nil)
	mockUserRepo.On("CreateUser", ctx, mock.AnythingOfType("*model.User")).Return(nil)

	resp, err := authService.Register(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.Email, resp.Email)
	assert.Equal(t, req.Name, resp.Name)
	mockUserRepo.AssertExpectations(t)
	mockCityRepo.AssertExpectations(t)
}

func TestAuthService_Register_UserAlreadyExists(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockCityRepo := new(MockCityRepository)
	mockCacheRepo := new(MockCacheRepository)
	jwtService := jwt.NewJWT(&jwt.JWTConfig{
		Secret:            "test-secret",
		TokenExpiration:   15 * time.Minute,
		RefreshExpiration: 7 * 24 * time.Hour,
	})

	authService := NewAuthService(mockUserRepo, mockCityRepo, mockCacheRepo, jwtService)

	ctx := context.Background()
	req := &contracts.RegisterRequest{
		Email:    "existing@example.com",
		Password: "password123",
		Name:     "John",
		Surname:  "Doe",
	}

	existingUser := &model.User{
		UserID: 1,
		Email:  req.Email,
	}

	mockUserRepo.On("GetUserByEmail", ctx, req.Email).Return(existingUser, nil)

	resp, err := authService.Register(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "already exists")
	mockUserRepo.AssertExpectations(t)
}

func TestAuthService_Login_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockCityRepo := new(MockCityRepository)
	mockCacheRepo := new(MockCacheRepository)
	jwtService := jwt.NewJWT(&jwt.JWTConfig{
		Secret:            "test-secret",
		TokenExpiration:   15 * time.Minute,
		RefreshExpiration: 7 * 24 * time.Hour,
	})

	authService := NewAuthService(mockUserRepo, mockCityRepo, mockCacheRepo, jwtService)

	ctx := context.Background()
	password := "password123"
	hashedPassword, _ := hash.Hash(password)

	user := &model.User{
		UserID:       1,
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
	}

	req := &contracts.LoginRequest{
		Email:    user.Email,
		Password: password,
	}

	mockUserRepo.On("GetUserByEmail", ctx, req.Email).Return(user, nil)

	resp, err := authService.Login(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken.Token)
	assert.NotEmpty(t, resp.RefreshToken.Token)
	assert.Greater(t, resp.AccessToken.ExpiresAt, int64(0))
	assert.Greater(t, resp.RefreshToken.ExpiresAt, int64(0))
	mockUserRepo.AssertExpectations(t)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockCityRepo := new(MockCityRepository)
	mockCacheRepo := new(MockCacheRepository)
	jwtService := jwt.NewJWT(&jwt.JWTConfig{
		Secret:            "test-secret",
		TokenExpiration:   15 * time.Minute,
		RefreshExpiration: 7 * 24 * time.Hour,
	})

	authService := NewAuthService(mockUserRepo, mockCityRepo, mockCacheRepo, jwtService)

	ctx := context.Background()
	req := &contracts.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	mockUserRepo.On("GetUserByEmail", ctx, req.Email).Return(nil, db.ErrNotFound)

	resp, err := authService.Login(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "does not exist")
	mockUserRepo.AssertExpectations(t)
}

func TestAuthService_Login_IncorrectPassword(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockCityRepo := new(MockCityRepository)
	mockCacheRepo := new(MockCacheRepository)
	jwtService := jwt.NewJWT(&jwt.JWTConfig{
		Secret:            "test-secret",
		TokenExpiration:   15 * time.Minute,
		RefreshExpiration: 7 * 24 * time.Hour,
	})

	authService := NewAuthService(mockUserRepo, mockCityRepo, mockCacheRepo, jwtService)

	ctx := context.Background()
	hashedPassword, _ := hash.Hash("correct_password")

	user := &model.User{
		UserID:       1,
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
	}

	req := &contracts.LoginRequest{
		Email:    user.Email,
		Password: "wrong_password",
	}

	mockUserRepo.On("GetUserByEmail", ctx, req.Email).Return(user, nil)

	resp, err := authService.Login(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "incorrect")
	mockUserRepo.AssertExpectations(t)
}

func TestAuthService_Logout_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockCityRepo := new(MockCityRepository)
	mockCacheRepo := new(MockCacheRepository)
	jwtService := jwt.NewJWT(&jwt.JWTConfig{
		Secret:            "test-secret",
		TokenExpiration:   15 * time.Minute,
		RefreshExpiration: 7 * 24 * time.Hour,
	})

	authService := NewAuthService(mockUserRepo, mockCityRepo, mockCacheRepo, jwtService)

	ctx := context.Background()
	accessToken, _ := jwtService.GenerateAccessToken(1, "test@example.com")
	refreshToken, _ := jwtService.GenerateRefreshToken(1, "test@example.com")

	req := &contracts.LogoutRequest{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	mockCacheRepo.On("SetBlacklist", ctx, accessToken, mock.AnythingOfType("time.Duration")).Return(nil)
	mockCacheRepo.On("SetBlacklist", ctx, refreshToken, mock.AnythingOfType("time.Duration")).Return(nil)

	err := authService.Logout(ctx, req)

	assert.NoError(t, err)
	mockCacheRepo.AssertExpectations(t)
}

func TestAuthService_RefreshToken_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockCityRepo := new(MockCityRepository)
	mockCacheRepo := new(MockCacheRepository)
	jwtService := jwt.NewJWT(&jwt.JWTConfig{
		Secret:            "test-secret",
		TokenExpiration:   15 * time.Minute,
		RefreshExpiration: 7 * 24 * time.Hour,
	})

	authService := NewAuthService(mockUserRepo, mockCityRepo, mockCacheRepo, jwtService)

	ctx := context.Background()
	refreshToken, _ := jwtService.GenerateRefreshToken(1, "test@example.com")

	req := &contracts.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	mockCacheRepo.On("IsBlacklisted", ctx, refreshToken).Return(false, nil)
	mockCacheRepo.On("SetBlacklist", ctx, refreshToken, mock.AnythingOfType("time.Duration")).Return(nil)

	resp, err := authService.RefreshToken(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken.Token)
	assert.NotEmpty(t, resp.RefreshToken.Token)
	mockCacheRepo.AssertExpectations(t)
}

func TestAuthService_RefreshToken_TokenBlacklisted(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockCityRepo := new(MockCityRepository)
	mockCacheRepo := new(MockCacheRepository)
	jwtService := jwt.NewJWT(&jwt.JWTConfig{
		Secret:            "test-secret",
		TokenExpiration:   15 * time.Minute,
		RefreshExpiration: 7 * 24 * time.Hour,
	})

	authService := NewAuthService(mockUserRepo, mockCityRepo, mockCacheRepo, jwtService)

	ctx := context.Background()
	refreshToken, _ := jwtService.GenerateRefreshToken(1, "test@example.com")

	req := &contracts.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	mockCacheRepo.On("IsBlacklisted", ctx, refreshToken).Return(true, nil)

	resp, err := authService.RefreshToken(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "revoked")
	mockCacheRepo.AssertExpectations(t)
}

func TestAuthService_RefreshToken_InvalidToken(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockCityRepo := new(MockCityRepository)
	mockCacheRepo := new(MockCacheRepository)
	jwtService := jwt.NewJWT(&jwt.JWTConfig{
		Secret:            "test-secret",
		TokenExpiration:   15 * time.Minute,
		RefreshExpiration: 7 * 24 * time.Hour,
	})

	authService := NewAuthService(mockUserRepo, mockCityRepo, mockCacheRepo, jwtService)

	ctx := context.Background()
	req := &contracts.RefreshTokenRequest{
		RefreshToken: "invalid-token",
	}

	resp, err := authService.RefreshToken(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}
