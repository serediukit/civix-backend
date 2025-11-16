package config

import (
	"os"

	"github.com/serediukit/civix-backend/pkg/database"
	"github.com/serediukit/civix-backend/pkg/jwt"
	"github.com/serediukit/civix-backend/pkg/redis"
)

type Config struct {
	Server   *ServerConfig
	Database *database.DatabaseConfig
	Redis    *redis.RedisConfig
	JWT      *jwt.JWTConfig
}

type ServerConfig struct {
	Port    string
	GinMode string
}

func LoadConfig() (*Config, error) {
	_ = os.Setenv("ENV_FILE_LOADED", "true")

	databaseConfig := database.GetDBConfig()
	redisConfig := redis.GetRedisConfig()
	jwtConfig := jwt.GetJWTConfig()

	return &Config{
		Server: &ServerConfig{
			Port:    getEnv("PORT", "8443"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		Database: databaseConfig,
		Redis:    redisConfig,
		JWT:      jwtConfig,
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

func (c *Config) GetJWTConfig() *jwt.JWTConfig {
	return jwt.GetJWTConfig()
}
