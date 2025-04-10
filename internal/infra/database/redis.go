package database

import (
	"context"
	"encoding/json"

	"log"

	databasedomain "github.com/israelalvesmelo/desafio-rate-limiter/internal/domain/database"
	"github.com/israelalvesmelo/desafio-rate-limiter/internal/domain/entity"

	"github.com/redis/go-redis/v9"
)

type RedisDb struct {
	client *redis.Client
}

func NewRedisDb(client *redis.Client) databasedomain.StorageDb {
	return &RedisDb{client: client}
}

func (r *RedisDb) SaveRateLimitConfig(ctx context.Context, apiKey *entity.RateLimitConfig) error {
	jsonReq, err := json.Marshal(apiKey)
	if err != nil {
		log.Println("error marshaling rate limit config")
		return err
	}

	if redisErr := r.client.Set(
		ctx,
		apiKey.Value,
		jsonReq,
		0,
	).Err(); redisErr != nil {
		log.Println("error inserting rate limit config")
		return redisErr
	}

	return nil
}
