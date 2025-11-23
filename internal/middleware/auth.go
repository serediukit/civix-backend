package middleware

import (
	"context"
	"errors"
	"strings"

	"github.com/serediukit/civix-backend/pkg/util/response"

	"github.com/gin-gonic/gin"
	"github.com/serediukit/civix-backend/internal/repository"
	"github.com/serediukit/civix-backend/pkg/jwt"
)

type AuthMiddleware struct {
	jwtUtil   *jwt.JWT
	redisRepo repository.CacheRepository
}

func NewAuthMiddleware(jwtUtil *jwt.JWT, redisRepo repository.CacheRepository) *AuthMiddleware {
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
			response.Unauthorized(c, "Authorization header is required", errors.New("missing authorization header"))
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "Invalid authorization header format", errors.New("invalid token format"))
			c.Abort()
			return
		}

		tokenString := parts[1]

		blacklisted, err := m.redisRepo.IsBlacklisted(c.Request.Context(), tokenString)
		if err != nil {
			response.InternalServerError(c, "Failed to verify token", err)
			c.Abort()
			return
		}

		if blacklisted {
			response.Unauthorized(c, "Token has been revoked", errors.New("token revoked"))
			c.Abort()
			return
		}

		claims, err := m.jwtUtil.ValidateToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "Invalid or expired token", err)
			c.Abort()
			return
		}

		ctx := context.WithValue(c.Request.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "user_email", claims.Email)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func GetUserIDFromContext(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value("user_id").(uint)
	return userID, ok
}

func GetUserEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value("user_email").(string)
	return email, ok
}
