package cache

import (
	"example-wikipedia-scraper/internal/config"
)

var redisClient *RedisClient

func InitRedis(cfg *config.RedisConfig) error {
	client, err := NewRedisClient(cfg)
	if err != nil {
		return err
	}
	redisClient = client
	return nil
}

func GetRedisClient() *RedisClient {
	return redisClient
}
