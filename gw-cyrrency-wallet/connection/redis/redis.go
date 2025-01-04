package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Stop()
}

type redisClient struct {
	conf string
	conn *redis.Client
}

func New(config string) RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	return &redisClient{conf: config, conn: client}
}

func (c *redisClient) Start() error {
	if err := c.conn.Ping(context.Background()).Err(); err != nil {
		return err
	}

	return nil
}

func (c *redisClient) Stop() {
	c.conn.Close()
}

func (r *redisClient) Get(ctx context.Context, key string) (string, error) {
	return r.conn.Get(ctx, key).Result()
}

func (r *redisClient) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return r.conn.Set(ctx, key, value, expiration).Err()
}
