package redisClient

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(options *redis.Options) *RedisClient {
	client := redis.NewClient(options)
	err := client.Conn().Ping(context.Background()).Err()
	if err != nil {
		log.Fatal("❌ can't connect to redis", err)
	}

	return &RedisClient{client: client}
}

func (rc *RedisClient) Set(key string, value string, ttl string) error {
	parsedTtl, err := time.ParseDuration(ttl)
	if err != nil {
		return err
	}

	return rc.client.Set(context.Background(), key, value, parsedTtl).Err()
}

func (rc *RedisClient) Get(key string) (string, error) {
	result, err := rc.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}

	return result, nil
}

func (rc *RedisClient) Delete(key string) error {
	return rc.client.Del(context.Background(), key).Err()
}
