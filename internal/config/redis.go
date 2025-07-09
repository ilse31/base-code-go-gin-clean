package config

import "fmt"

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// RedisURL returns the Redis connection URL
func (r *RedisConfig) RedisURL() string {
	if r.Password != "" {
		return fmt.Sprintf("redis://:%s@%s:%s/%d", r.Password, r.Host, r.Port, r.DB)
	}
	return fmt.Sprintf("redis://%s:%s/%d", r.Host, r.Port, r.DB)
}
