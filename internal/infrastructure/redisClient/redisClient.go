package redisClient

import "github.com/redis/go-redis/v9"

func NewRedisClient(options *redis.Options) *redis.Client {
	client := redis.NewClient(options)

	return client
}
