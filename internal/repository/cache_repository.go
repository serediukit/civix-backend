package repository

import (
	"context"
	"time"

	"github.com/serediukit/civix-backend/pkg/redis"
)

type CacheRepository interface {
	SetBlacklist(ctx context.Context, token string, expiration time.Duration) error
	IsBlacklisted(ctx context.Context, token string) (bool, error)
	Delete(ctx context.Context, key string) error
}

type cacheRepository struct {
	cachedStore *redis.CachedStore
}

func NewCacheRepository(cachedStore *redis.CachedStore) CacheRepository {
	return &cacheRepository{cachedStore: cachedStore}
}

func (r *cacheRepository) SetBlacklist(ctx context.Context, token string, expiration time.Duration) error {
	return r.cachedStore.GetClient().Set(ctx, "blacklist:"+token, "1", expiration).Err()
}

func (r *cacheRepository) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	val, err := r.cachedStore.GetClient().Exists(ctx, "blacklist:"+token).Result()
	if err != nil {
		return false, err
	}
	return val > 0, nil
}

func (r *cacheRepository) Delete(ctx context.Context, key string) error {
	return r.cachedStore.GetClient().Del(ctx, key).Err()
}
