package database

import (
	"context"

	"github.com/israelalvesmelo/desafio-rate-limiter/internal/domain/entity"
)

type StorageDb interface {
	SaveRateLimitConfig(ctx context.Context, apiKey *entity.RateLimitConfig) error
}
