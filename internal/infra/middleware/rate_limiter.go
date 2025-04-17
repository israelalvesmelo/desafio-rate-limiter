package middleware

import (
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
	config  config.Config
}

func NewRateLimitMiddleware(storage database.StorageDb, config config.Config) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		storage: storage,
		config:  config,
	}
}

func (m *RateLimiterMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := getKey(r)
		//TODO: Aqui tenho que buscar o os valores da confg no redis ou no banco de dados
		limiter := usecase.NewLimiter(m.storage, m.config)

		resp, err := limiter.Execute(r.Context(),
			dto.RequestSave{
				Key:      key,
				TimeAdd: time.Now(),
			})
		if errors.Is(err, entity.ErrIPExceededAmountRequest) {
			log.Printf("Error executing request: %s\n", err.Error())
			http.Error(w, err.Error(), http.StatusTooManyRequests)
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

func getKey(r *http.Request) string {
	apiKey := r.Header.Get(entity.APIKeyHeaderName)
	if apiKey != "" {
		return apiKey
	}
	ip := getClientIP(r)
	return ip
}

// getClientIP extracts the client IP from the request
func getClientIP(r *http.Request) string {
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
