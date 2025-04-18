package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/israelalvesmelo/desafio-rate-limiter/cmd/api/config"
	"github.com/israelalvesmelo/desafio-rate-limiter/internal/domain/database"
	"github.com/israelalvesmelo/desafio-rate-limiter/internal/domain/entity"
	"github.com/israelalvesmelo/desafio-rate-limiter/internal/domain/usecase"
	"github.com/israelalvesmelo/desafio-rate-limiter/internal/infra"
	"github.com/israelalvesmelo/desafio-rate-limiter/internal/infra/dto"
	"github.com/israelalvesmelo/desafio-rate-limiter/internal/infra/handler"
)

type RateLimiterMiddleware struct {
	storage database.StorageDb
	config  config.RateLimiter
}

func NewRateLimitMiddleware(storage database.StorageDb, config config.RateLimiter) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		storage: storage,
		config:  config,
	}
}

func (m *RateLimiterMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key, keyType := m.getKey(r)
		rtConfig := m.getRateLimitConfig(r.Context(), key, keyType)

		limiter := usecase.NewLimiter(m.storage)
		resp, err := limiter.Execute(r.Context(),
			dto.RequestSave{
				TimeAdd: time.Now(),
			},
			rtConfig,
		)
		if errors.Is(err, entity.ErrExceededAmountRequest) {
			log.Printf("Error executing request: %s\n", err.Error())
			handler.Error(w, err.Error(), http.StatusTooManyRequests)
			return
		}
		if err != nil {
			log.Printf("Error executing request: %s\n", err.Error())
			handler.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !resp.Allow {
			log.Printf("Too many request: %s\n", entity.ErrExceededAmountRequest.Error())
			handler.Error(w, entity.ErrExceededAmountRequest.Error(), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *RateLimiterMiddleware) getKey(r *http.Request) (string, string) {
	apiKey := r.Header.Get(entity.APIKeyName)
	if apiKey != "" {
		return apiKey, entity.APIKeyName
	}
	ip := infra.GetClientIP(r)
	return ip, entity.IPName
}

func (m *RateLimiterMiddleware) getRateLimitConfig(ctx context.Context, key string, keyType string) entity.RateLimitConfig {
	rlConfig, _ := m.storage.GetRateLimitConfig(ctx, key)
	if rlConfig != nil {
		return *rlConfig
	}

	limitValues := entity.LimitValues(m.config.ByAPIKey)
	if keyType == entity.IPName {
		limitValues = entity.LimitValues(m.config.ByIP)
	}

	return entity.RateLimitConfig{
		Key:         key,
		LimitValues: limitValues,
	}
}
