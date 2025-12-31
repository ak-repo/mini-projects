package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type PresenceRepository struct {
	client *redis.Client
	ctx    context.Context
}

func NewPresenceRepository(client *redis.Client) *PresenceRepository {
	return &PresenceRepository{
		client: client,
		ctx:    context.Background(),
	}
}

func (r *PresenceRepository) SetUserOnline(userID string, ttl time.Duration) error {
	key := fmt.Sprintf("presence:user:%s", userID)
	return r.client.Set(r.ctx, key, "online", ttl).Err()
}

func (r *PresenceRepository) SetUserOffline(userID string) error {
	key := fmt.Sprintf("presence:user:%s", userID)
	return r.client.Del(r.ctx, key).Err()
}

func (r *PresenceRepository) IsUserOnline(userID string) (bool, error) {
	key := fmt.Sprintf("presence:user:%s", userID)
	val, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return val == "online", nil
}

func (r *PresenceRepository) SetLastSeen(userID string) error {
	key := fmt.Sprintf("last_seen:user:%s", userID)
	timestamp := time.Now().Unix()
	return r.client.Set(r.ctx, key, timestamp, 0).Err()
}

func (r *PresenceRepository) GetLastSeen(userID string) (*time.Time, error) {
	key := fmt.Sprintf("last_seen:user:%s", userID)
	val, err := r.client.Get(r.ctx, key).Int64()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	t := time.Unix(val, 0)
	return &t, nil
}
