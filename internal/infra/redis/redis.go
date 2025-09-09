package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	*redis.Client
}

func NewClient(host string, port int, password string, db int) (*Client, func(), error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	cleanup := func() {
		_ = rdb.Close()
	}

	return &Client{Client: rdb}, cleanup, nil
}

func (c *Client) SetLock(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	result := c.SetNX(ctx, key, value, expiration)
	if result.Err() != nil {
		return false, fmt.Errorf("failed to set lock: %w", result.Err())
	}
	return result.Val(), nil
}

func (c *Client) ReleaseLock(ctx context.Context, key string) error {
	result := c.Del(ctx, key)
	if result.Err() != nil {
		return fmt.Errorf("failed to release lock: %w", result.Err())
	}
	return nil
}

func (c *Client) Close() error {
	return c.Client.Close()
}
