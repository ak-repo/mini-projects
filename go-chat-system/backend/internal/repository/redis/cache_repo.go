package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheRepository struct {
	client *redis.Client
	ctx    context.Context
}

func NewCacheRepository(client *redis.Client) *CacheRepository {
	return &CacheRepository{
		client: client,
		ctx:    context.Background(),
	}
}

func (r *CacheRepository) Set(key string, value interface{}, ttl time.Duration) error {
	return r.client.Set(r.ctx, key, value, ttl).Err()
}

func (r *CacheRepository) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

func (r *CacheRepository) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}
