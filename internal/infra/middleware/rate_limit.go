package middleware

import (
	"log"
	"net/http"
)

type RateLimitMiddleware struct {
}

func NewRateLimitMiddleware() *RateLimitMiddleware {
	return &RateLimitMiddleware{}
}

func (m *RateLimitMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Implemente a lógica de rate limit aqui
		log.Println("Rate limit middleware...")
		next.ServeHTTP(w, r)
	})
}
