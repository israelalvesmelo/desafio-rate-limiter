package database

import (
	"context"

	"github.com/israelalvesmelo/desafio-rate-limiter/internal/domain/entity"
)

type StorageDb interface {
	// SaveRateLimitConfig Stores the rate limit configuration customized by the API key or IP.
	SaveRateLimitConfig(ctx context.Context, Key *entity.RateLimitConfig) error

	// GetRateLimitConfig Obtains the rate limit configuration customized by the API key or IP.
	GetRateLimitConfig(ctx context.Context, key string) (*entity.RateLimitConfig, error)

	// UpsertRequest Updates or inserts a new request inside the RateLimiter.Requests array.
	// This also creates a new instance of Request
	UpsertRequest(ctx context.Context, key string, rl *entity.RateLimiter) error

	// SaveBlockedDuration Stores the blocked duration amount by key
	SaveBlockedDuration(ctx context.Context, key string, BlockedDuration int64) error

	// GetBlockedDuration Obtain the blocked duration by key
	GetBlockedDuration(ctx context.Context, key string) (string, error)

	// GetRequest reads the stored array of request
	GetRequest(ctx context.Context, key string) (*entity.RateLimiter, error)
}
