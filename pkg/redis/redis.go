package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var retriesCount = 5

type CachedStore struct {
	client *redis.Client
}

func (s *CachedStore) Close() {
	if err := s.client.Close(); err != nil {
		log.Fatalf("[Redis] Failed to close connection: %v", err)
	}

	log.Println("[Redis] connection closed")
}

func NewRedis(config *RedisConfig) (*CachedStore, error) {
	var err error

	for retriesCount > 0 {
		client := redis.NewClient(&redis.Options{
			Addr:         fmt.Sprintf("%s:%s", config.Host, config.Port),
			Password:     config.Password,
			DB:           config.DB,
			PoolSize:     config.PoolSize,
			MinIdleConns: config.MinIdleConns,
			DialTimeout:  config.DialTimeout,
			ReadTimeout:  config.ReadTimeout,
		})

		_, err = client.Ping(context.Background()).Result()
		if err == nil {
			return &CachedStore{
				client: client,
			}, nil
		}
		log.Println("[Redis] Waiting for ...")
		time.Sleep(2 * time.Second)

		retriesCount--
	}

	log.Fatalf("[Redis] Unable to ping: %v\n", err)
	return nil, err
}

func (s *CachedStore) GetClient() *redis.Client {
	return s.client
}
