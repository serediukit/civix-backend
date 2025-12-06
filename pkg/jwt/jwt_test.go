package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestNewJWT(t *testing.T) {
	config := &JWTConfig{
		Secret:            "test-secret",
		TokenExpiration:   15 * time.Minute,
		RefreshExpiration: 7 * 24 * time.Hour,
	}

	jwtService := NewJWT(config)

	assert.NotNil(t, jwtService)
	assert.Equal(t, []byte(config.Secret), jwtService.secretKey)
	assert.Equal(t, config.TokenExpiration, jwtService.accessTokenExpiration)
	assert.Equal(t, config.RefreshExpiration, jwtService.refreshTokenExpiration)
}

func TestJWT_GenerateAccessToken_Success(t *testing.T) {
	jwtService := NewJWT(&JWTConfig{
		Secret:            "test-secret",
		TokenExpiration:   15 * time.Minute,
		RefreshExpiration: 7 * 24 * time.Hour,
	})

	userID := uint64(123)
	email := "test@example.com"

	token, err := jwtService.GenerateAccessToken(userID, email)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify token can be parsed
	claims, err := jwtService.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
}

func TestJWT_GenerateRefreshToken_Success(t *testing.T) {
	jwtService := NewJWT(&JWTConfig{
		Secret:            "test-secret",
		TokenExpiration:   15 * time.Minute,
		RefreshExpiration: 7 * 24 * time.Hour,
	})

	userID := uint64(456)
	email := "refresh@example.com"

	token, err := jwtService.GenerateRefreshToken(userID, email)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify token can be parsed
	claims, err := jwtService.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
}

func TestJWT_GenerateToken_ClaimsPopulated(t *testing.T) {
	jwtService := NewJWT(&JWTConfig{
		Secret:            "test-secret",
		TokenExpiration:   15 * time.Minute,
		RefreshExpiration: 7 * 24 * time.Hour,
	})

	userID := uint64(789)
	email := "claims@example.com"

	token, err := jwtService.GenerateAccessToken(userID, email)
	assert.NoError(t, err)

	claims, err := jwtService.ValidateToken(token)
	assert.NoError(t, err)

	// Verify all claims are populated
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.NotNil(t, claims.ExpiresAt)
	assert.NotNil(t, claims.IssuedAt)
	assert.True(t, claims.ExpiresAt.After(time.Now()))
}

func TestJWT_ValidateToken_Success(t *testing.T) {
	jwtService := NewJWT(&JWTConfig{
		Secret:            "test-secret",
		TokenExpiration:   15 * time.Minute,
		RefreshExpiration: 7 * 24 * time.Hour,
	})

	userID := uint64(111)
	email := "valid@example.com"

	token, err := jwtService.GenerateAccessToken(userID, email)
	assert.NoError(t, err)

	claims, err := jwtService.ValidateToken(token)

	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
}

func TestJWT_ValidateToken_InvalidToken(t *testing.T) {
	jwtService := NewJWT(&JWTConfig{
		Secret:            "test-secret",
		TokenExpiration:   15 * time.Minute,
		RefreshExpiration: 7 * 24 * time.Hour,
	})

	invalidToken := "invalid.token.string"

	claims, err := jwtService.ValidateToken(invalidToken)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestJWT_ValidateToken_WrongSecret(t *testing.T) {
	jwtService1 := NewJWT(&JWTConfig{
		Secret:            "secret-1",
		TokenExpiration:   15 * time.Minute,
		RefreshExpiration: 7 * 24 * time.Hour,
	})

	jwtService2 := NewJWT(&JWTConfig{
		Secret:            "secret-2",
		TokenExpiration:   15 * time.Minute,
		RefreshExpiration: 7 * 24 * time.Hour,
	})

	token, err := jwtService1.GenerateAccessToken(1, "test@example.com")
	assert.NoError(t, err)

	// Try to validate with wrong secret
	claims, err := jwtService2.ValidateToken(token)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestJWT_ValidateToken_ExpiredToken(t *testing.T) {
	jwtService := NewJWT(&JWTConfig{
		Secret:            "test-secret",
		TokenExpiration:   -1 * time.Hour, // Already expired
		RefreshExpiration: 7 * 24 * time.Hour,
	})

	token, err := jwtService.GenerateAccessToken(1, "test@example.com")
	assert.NoError(t, err)

	claims, err := jwtService.ValidateToken(token)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestJWT_ValidateToken_InvalidSigningMethod(t *testing.T) {
	jwtService := NewJWT(&JWTConfig{
		Secret:            "test-secret",
		TokenExpiration:   15 * time.Minute,
		RefreshExpiration: 7 * 24 * time.Hour,
	})

	// Create a token with a different signing method (RS256 instead of HS256)
	claims := &Claims{
		UserID: 1,
		Email:  "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Using None signing method (invalid)
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenString, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)

	validatedClaims, err := jwtService.ValidateToken(tokenString)

	assert.Error(t, err)
	assert.Nil(t, validatedClaims)
	assert.Contains(t, err.Error(), "invalid token")
}

func TestJWT_ValidateToken_EmptyToken(t *testing.T) {
	jwtService := NewJWT(&JWTConfig{
		Secret:            "test-secret",
		TokenExpiration:   15 * time.Minute,
		RefreshExpiration: 7 * 24 * time.Hour,
	})

	claims, err := jwtService.ValidateToken("")

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestGenerateRandomString_Success(t *testing.T) {
	length := 16

	randomString, err := GenerateRandomString(length)

	assert.NoError(t, err)
	assert.NotEmpty(t, randomString)
	assert.Equal(t, length, len(randomString))
}

func TestGenerateRandomString_DifferentOutputs(t *testing.T) {
	length := 32

	string1, err1 := GenerateRandomString(length)
	string2, err2 := GenerateRandomString(length)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEmpty(t, string1)
	assert.NotEmpty(t, string2)
	assert.NotEqual(t, string1, string2)
}

func TestGenerateRandomString_ZeroLength(t *testing.T) {
	randomString, err := GenerateRandomString(0)

	assert.NoError(t, err)
	assert.Empty(t, randomString)
}

func TestGenerateRandomString_LargeLength(t *testing.T) {
	length := 1024

	randomString, err := GenerateRandomString(length)

	assert.NoError(t, err)
	assert.NotEmpty(t, randomString)
	assert.Equal(t, length, len(randomString))
}

func TestJWT_TokenExpirationTimes(t *testing.T) {
	accessExpiration := 15 * time.Minute
	refreshExpiration := 7 * 24 * time.Hour

	jwtService := NewJWT(&JWTConfig{
		Secret:            "test-secret",
		TokenExpiration:   accessExpiration,
		RefreshExpiration: refreshExpiration,
	})

	userID := uint64(1)
	email := "test@example.com"

	// Test access token expiration
	accessToken, err := jwtService.GenerateAccessToken(userID, email)
	assert.NoError(t, err)

	accessClaims, err := jwtService.ValidateToken(accessToken)
	assert.NoError(t, err)

	expectedAccessExpiry := time.Now().Add(accessExpiration)
	actualAccessExpiry := accessClaims.ExpiresAt.Time

	// Allow 1 second tolerance
	assert.WithinDuration(t, expectedAccessExpiry, actualAccessExpiry, 1*time.Second)

	// Test refresh token expiration
	refreshToken, err := jwtService.GenerateRefreshToken(userID, email)
	assert.NoError(t, err)

	refreshClaims, err := jwtService.ValidateToken(refreshToken)
	assert.NoError(t, err)

	expectedRefreshExpiry := time.Now().Add(refreshExpiration)
	actualRefreshExpiry := refreshClaims.ExpiresAt.Time

	// Allow 1 second tolerance
	assert.WithinDuration(t, expectedRefreshExpiry, actualRefreshExpiry, 1*time.Second)
}

func TestJWT_MultipleUserTokens(t *testing.T) {
	jwtService := NewJWT(&JWTConfig{
		Secret:            "test-secret",
		TokenExpiration:   15 * time.Minute,
		RefreshExpiration: 7 * 24 * time.Hour,
	})

	users := []struct {
		userID uint64
		email  string
	}{
		{1, "user1@example.com"},
		{2, "user2@example.com"},
		{3, "user3@example.com"},
	}

	for _, user := range users {
		token, err := jwtService.GenerateAccessToken(user.userID, user.email)
		assert.NoError(t, err)

		claims, err := jwtService.ValidateToken(token)
		assert.NoError(t, err)
		assert.Equal(t, user.userID, claims.UserID)
		assert.Equal(t, user.email, claims.Email)
	}
}
