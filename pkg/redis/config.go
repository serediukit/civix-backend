package redis

import (
	"time"

	"github.com/serediukit/civix-backend/pkg/env"
)

type RedisConfig struct {
	Host         string
	Port         string
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
}

func GetRedisConfig() *RedisConfig {
	return &RedisConfig{
		Host:         env.GetEnv("REDIS_HOST", "localhost"),
		Port:         env.GetEnv("REDIS_PORT", "6379"),
		Password:     env.GetEnv("REDIS_PASSWORD", ""),
		DB:           env.GetEnvInt("REDIS_DB", 0),
		PoolSize:     env.GetEnvInt("REDIS_POOLSIZE", 10),
		MinIdleConns: env.GetEnvInt("REDIS_MIDIDLECONNS", 5),
		DialTimeout:  env.GetEnvDurationSeconds("REDIS_DB", 5*time.Second),
		ReadTimeout:  env.GetEnvDurationSeconds("REDIS_DB", 5*time.Second),
	}
}
