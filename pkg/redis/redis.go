package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type CachedStore struct {
	client *redis.Client
}

func (s CachedStore) Close() {
	if err := s.client.Close(); err != nil {
		log.Fatalf("Failed to close Redis connection: %v", err)
	}
}

func NewRedis(config *RedisConfig) (*CachedStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
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
