package middleware

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/israelalvesmelo/desafio-rate-limiter/cmd/api/config"
	"github.com/israelalvesmelo/desafio-rate-limiter/internal/domain/database"
	"github.com/israelalvesmelo/desafio-rate-limiter/internal/domain/entity"
	"github.com/israelalvesmelo/desafio-rate-limiter/internal/domain/usecase"
	"github.com/israelalvesmelo/desafio-rate-limiter/internal/infra/dto"
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
		if errors.Is(err, entity.ErrIPExceededAmountRequest) {
			log.Printf("Error executing request: %s\n", err.Error())
			http.Error(w, err.Error(), http.StatusTooManyRequests) //TODO: MUDAR PARA JSON
			return
		}
		if err != nil {
			log.Printf("Error executing request: %s\n", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !resp.Allow {
			log.Printf("Too many request: %s\n", entity.ErrIPExceededAmountRequest.Error())
			http.Error(w, entity.ErrIPExceededAmountRequest.Error(), http.StatusTooManyRequests)
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
	ip := m.getClientIP(r)
	return ip, entity.IPName
}

// getClientIP extracts the client IP from the request
func (m *RateLimiterMiddleware) getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := net.ParseIP(xff)
		if ips != nil {
			return ips.String()
		}
	}

	// Extract from RemoteAddr
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
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
