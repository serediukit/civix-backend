package jwt

import (
	"github.com/serediukit/civix-backend/pkg/util/timeutil"
	"time"

	"github.com/serediukit/civix-backend/pkg/env"
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
