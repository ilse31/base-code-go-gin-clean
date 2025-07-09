package redis

import (
	"context"
	"time"

	"base-code-go-gin-clean/internal/config"

	"github.com/redis/go-redis/v9"
)

// Client represents a Redis client wrapper
type Client struct {
	client *redis.Client
}

// NewClient creates a new Redis client
func NewClient(cfg *config.RedisConfig) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Client{client: client}, nil
}

// Close closes the Redis client connection
func (c *Client) Close() error {
	return c.client.Close()
}

// GetClient returns the underlying Redis client
func (c *Client) GetClient() *redis.Client {
	return c.client
}
