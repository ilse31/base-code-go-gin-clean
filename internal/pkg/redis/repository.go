package redis

import (
	"context"
	"time"
)

// Repository defines the interface for Redis operations
type Repository interface {
	// Set sets a key-value pair with an expiration time
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	// Get retrieves a value by key
	Get(ctx context.Context, key string) (string, error)
	// Delete removes a key
	Delete(ctx context.Context, key string) error
	// Exists checks if a key exists
	Exists(ctx context.Context, key string) (bool, error)
	// Expire sets a timeout on a key
	Expire(ctx context.Context, key string, expiration time.Duration) (bool, error)
	// Close closes the Redis connection
	Close() error
}

// redisRepository implements the Repository interface
type redisRepository struct {
	client *Client
}

// NewRepository creates a new Redis repository
func NewRepository(client *Client) Repository {
	return &redisRepository{client: client}
}

// Set sets a key-value pair with an expiration time
func (r *redisRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value by key
func (r *redisRepository) Get(ctx context.Context, key string) (string, error) {
	return r.client.client.Get(ctx, key).Result()
}

// Delete removes a key
func (r *redisRepository) Delete(ctx context.Context, key string) error {
	return r.client.client.Del(ctx, key).Err()
}

// Exists checks if a key exists
func (r *redisRepository) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := r.client.client.Exists(ctx, key).Result()
	return exists > 0, err
}

// Expire sets a timeout on a key
func (r *redisRepository) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return r.client.client.Expire(ctx, key, expiration).Result()
}

// Close closes the Redis connection
func (r *redisRepository) Close() error {
	return r.client.Close()
}
