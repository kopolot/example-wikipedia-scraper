package cache

import (
	"context"
	"fmt"
	"time"

	"example-wikipedia-scraper/internal/config"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisClient(cfg *config.RedisConfig) (*RedisClient, error) {
	if cfg == nil {
		return nil, fmt.Errorf("redis config is nil")
	}
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}
	ttl := time.Duration(cfg.DefaultTTLSeconds) * time.Second
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	return &RedisClient{client: client, ttl: ttl}, nil
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisClient) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	if ttl <= 0 {
		ttl = r.ttl
	}
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisClient) SetNX(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
	if ttl <= 0 {
		ttl = r.ttl
	}
	return r.client.SetNX(ctx, key, value, ttl).Result()
}

func (r *RedisClient) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}
