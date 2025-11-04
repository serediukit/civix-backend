package middleware

import (
	"context"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/serediukit/civix-backend/internal/repository"
	"github.com/serediukit/civix-backend/internal/util"
)

type AuthMiddleware struct {
	jwtUtil   *util.JWTUtil
	redisRepo repository.RedisRepository
}

func NewAuthMiddleware(jwtUtil *util.JWTUtil, redisRepo repository.RedisRepository) *AuthMiddleware {
	return &AuthMiddleware{
		jwtUtil:   jwtUtil,
		redisRepo: redisRepo,
	}
}

// AuthRequired is a middleware that checks for a valid JWT token in the Authorization header
func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			util.Unauthorized(c, "Authorization header is required", errors.New("missing authorization header"))
			c.Abort()
			return
		}

		// Check if the token starts with "Bearer "
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			util.Unauthorized(c, "Invalid authorization header format", errors.New("invalid token format"))
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Check if token is blacklisted
		blacklisted, err := m.redisRepo.IsBlacklisted(c.Request.Context(), tokenString)
		if err != nil {
			util.InternalServerError(c, "Failed to verify token", err)
			c.Abort()
			return
		}

		if blacklisted {
			util.Unauthorized(c, "Token has been revoked", errors.New("token revoked"))
			c.Abort()
			return
		}

		// Validate token
		claims, err := m.jwtUtil.ValidateToken(tokenString)
		if err != nil {
			util.Unauthorized(c, "Invalid or expired token", err)
			c.Abort()
			return
		}

		// Add user info to context
		ctx := context.WithValue(c.Request.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "user_email", claims.Email)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// GetUserIDFromContext retrieves the user UserID from the context
func GetUserIDFromContext(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value("user_id").(uint)
	return userID, ok
}

// GetUserEmailFromContext retrieves the user email from the context
func GetUserEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value("user_email").(string)
	return email, ok
}
