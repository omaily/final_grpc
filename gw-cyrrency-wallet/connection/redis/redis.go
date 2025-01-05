package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/omaily/final_grpc/gw-cyrrency-wallet/config"
)

type Client interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value []byte, expiration time.Duration) error
	Stop()
}

type redisClient struct {
	conf config.RedisServer
	conn *redis.Client
}

func New(conf config.RedisServer) Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Address, conf.Port),
		Password: "",
		DB:       0,
	})

	return &redisClient{conf: conf, conn: client}
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

func (r *redisClient) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	return r.conn.Set(ctx, key, value, expiration).Err()
}
