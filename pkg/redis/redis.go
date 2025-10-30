package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisConfigGetter interface {
	GetRedisConfig() *RedisConfig
}

type CachedStore struct {
	client *redis.Client
}

func NewRedis(config RedisConfigGetter) (*CachedStore, error) {
	redisConfig := config.GetRedisConfig()

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", redisConfig.Host, redisConfig.Port),
		Password:     redisConfig.Password,
		DB:           redisConfig.DB,
		PoolSize:     redisConfig.PoolSize,
		MinIdleConns: redisConfig.MinIdleConns,
		DialTimeout:  redisConfig.DialTimeout,
		ReadTimeout:  redisConfig.ReadTimeout,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("[Redis] Unable to ping: %v\n", err)
		return nil, err
	}

	return &CachedStore{
		client: client,
	}, nil
}
