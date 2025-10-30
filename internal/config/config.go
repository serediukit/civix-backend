package config

import (
	"os"
	"time"

	"github.com/serediukit/civix-backend/pkg/database"
	"github.com/serediukit/civix-backend/pkg/redis"
)

type Config struct {
	Server   *ServerConfig
	Database *database.DatabaseConfig
	Redis    *redis.RedisConfig
	JWT      *JWTConfig
}

type ServerConfig struct {
	Port    string
	GinMode string
}

type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

func LoadConfig() (*Config, error) {
	_ = os.Setenv("ENV_FILE_LOADED", "true")

	redisConfig := redis.GetRedisConfig()
	databaseConfig := database.GetDBConfig()

	jwtExpiration, _ := time.ParseDuration(getEnv("JWT_EXPIRATION", "24h"))

	return &Config{
		Server: &ServerConfig{
			Port:    getEnv("PORT", "8080"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		Database: databaseConfig,
		Redis:    redisConfig,
		JWT: &JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your_jwt_secret_key_here"),
			Expiration: jwtExpiration,
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func (c *Config) GetDBConfig() *database.DatabaseConfig {
	return c.Database
}

func (c *Config) GetRedisConfig() *redis.RedisConfig {
	return c.Redis
}
