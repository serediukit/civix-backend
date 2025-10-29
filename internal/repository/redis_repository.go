package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository interface {
	SetBlacklist(ctx context.Context, token string, expiration time.Duration) error
	IsBlacklisted(ctx context.Context, token string) (bool, error)
	Delete(ctx context.Context, key string) error
}

type redisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) RedisRepository {
	return &redisRepository{client: client}
}

func (r *redisRepository) SetBlacklist(ctx context.Context, token string, expiration time.Duration) error {
	return r.client.Set(ctx, "blacklist:"+token, "1", expiration).Err()
}

func (r *redisRepository) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	val, err := r.client.Exists(ctx, "blacklist:"+token).Result()
	if err != nil {
		return false, err
	}
	return val > 0, nil
}

func (r *redisRepository) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
