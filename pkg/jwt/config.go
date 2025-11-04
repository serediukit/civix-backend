package jwt

import (
	"time"

	"github.com/serediukit/civix-backend/pkg/env"
	"github.com/serediukit/civix-backend/pkg/timeutil"
)

type JWTConfig struct {
	Secret            string
	TokenExpiration   time.Duration
	RefreshExpiration time.Duration
}

func GetJWTConfig() *JWTConfig {
	return &JWTConfig{
		Secret:            env.GetEnv("JWT_SECRET", "your_jwt_secret_key_here"),
		TokenExpiration:   env.GetEnvDurationSeconds("JWT_TOKEN_EXPIRATION", 15*time.Minute),
		RefreshExpiration: env.GetEnvDurationSeconds("JWT_REFRESH_EXPIRATION", timeutil.Week),
	}
}
