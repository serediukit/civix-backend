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

	log.Println(fmt.Sprintf("%s:%s", config.Host, config.Port))
	log.Println(config.Password)
	log.Println(config.DB)
	log.Println(config.PoolSize)
	log.Println(config.MinIdleConns)
	log.Println(config.DialTimeout)
	log.Println(config.ReadTimeout)

	var err error

	for retriesCount > 0 {
		_, err = client.Ping(context.Background()).Result()
		if err == nil {
			return &CachedStore{
				client: client,
			}, nil
		}
		log.Println("Waiting for Redis...")
		time.Sleep(2 * time.Second)

		retriesCount--
	}

	log.Fatalf("[Redis] Unable to ping: %v\n", err)
	return nil, err
}

func (s *CachedStore) GetClient() *redis.Client {
	return s.client
}
