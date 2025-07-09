package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"base-code-go-gin-clean/internal/domain/user"
	"base-code-go-gin-clean/internal/pkg/redis"
	"base-code-go-gin-clean/internal/pkg/telemetry"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type UserService interface {
	GetUserByID(ctx context.Context, id string) (*user.UserResponse, error)
}

type userService struct {
	userRepo  user.UserRepository
	redisRepo redis.Repository
	cacheTTL  time.Duration
}

// Cache key prefixes
const (
	userCacheKeyPrefix = "user:"
)

// Cache TTLs
const (
	defaultCacheTTL = 5 * time.Minute
)

type UserServiceConfig struct {
	UserRepo  user.UserRepository
	RedisRepo redis.Repository
	CacheTTL  time.Duration
}

func NewUserService(cfg UserServiceConfig) UserService {
	svc := &userService{
		userRepo:  cfg.UserRepo,
		redisRepo: cfg.RedisRepo,
		cacheTTL:  defaultCacheTTL,
	}

	// Override default cache TTL if provided
	if cfg.CacheTTL > 0 {
		svc.cacheTTL = cfg.CacheTTL
	}

	return svc
}

func (s *userService) getUserCacheKey(id string) string {
	return userCacheKeyPrefix + id
}

// GetUserByID retrieves a user by ID with caching
func (s *userService) GetUserByID(ctx context.Context, idStr string) (*user.UserResponse, error) {
	// Start a new span for the service method
	ctx, span := telemetry.StartSpan(ctx, "UserService.GetUserByID")
	defer telemetry.EndSpan(span, nil)

	// Add user ID to span
	span.SetAttributes(attribute.String("user.id", idStr))

	// Validate UUID format
	if _, err := uuid.Parse(idStr); err != nil {
		err = fmt.Errorf("invalid user ID format: %v", err)
		telemetry.EndSpan(span, err)
		return nil, err
	}

	// Try to get from cache first
	cacheKey := s.getUserCacheKey(idStr)
	cachedUser, err := s.getUserFromCache(ctx, cacheKey)
	if err == nil && cachedUser != nil {
		span.AddEvent("Cache hit", trace.WithAttributes(
			attribute.String("cache.key", cacheKey),
		))
		return cachedUser, nil
	}

	span.AddEvent("Cache miss", trace.WithAttributes(
		attribute.String("cache.key", cacheKey),
	))

	// Parse UUID for repository call
	userID, err := uuid.Parse(idStr)
	if err != nil {
		err = fmt.Errorf("invalid user ID format: %v", err)
		telemetry.EndSpan(span, err)
		return nil, err
	}

	// If not in cache, get from repository
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		err = fmt.Errorf("failed to get user from repository: %v", err)
		telemetry.EndSpan(span, err)
		return nil, err
	}

	// Convert to response model
	userResponse := user.ToResponse()

	// Cache the result
	if err := s.cacheUser(ctx, cacheKey, userResponse); err != nil {
		// Add cache error to span but don't fail the request
		errAttrs := trace.WithAttributes(
			attribute.String("cache.error", err.Error()),
		)
		span.RecordError(err, errAttrs)
	}

	return userResponse, nil
}

// getUserFromCache retrieves a user from the cache
func (s *userService) getUserFromCache(ctx context.Context, key string) (*user.UserResponse, error) {
	_, span := telemetry.StartSpan(ctx, "UserService.getUserFromCache")
	defer telemetry.EndSpan(span, nil)

	data, err := s.redisRepo.Get(ctx, key)
	if err != nil {
		telemetry.EndSpan(span, err)
		return nil, err
	}

	var userResp user.UserResponse
	if err := json.Unmarshal([]byte(data), &userResp); err != nil {
		err = fmt.Errorf("failed to unmarshal cached user: %v", err)
		telemetry.EndSpan(span, err)
		return nil, err
	}

	return &userResp, nil
}

// cacheUser stores a user in the cache
func (s *userService) cacheUser(ctx context.Context, key string, user *user.UserResponse) error {
	_, span := telemetry.StartSpan(ctx, "UserService.cacheUser")
	defer telemetry.EndSpan(span, nil)

	data, err := json.Marshal(user)
	if err != nil {
		err = fmt.Errorf("failed to marshal user for caching: %v", err)
		telemetry.EndSpan(span, err)
		return err
	}

	err = s.redisRepo.Set(ctx, key, string(data), s.cacheTTL)
	if err != nil {
		telemetry.EndSpan(span, err)
	} else {
		span.SetAttributes(
			attribute.String("cache.key", key),
			attribute.String("cache.ttl", s.cacheTTL.String()),
		)
	}

	return err
}
