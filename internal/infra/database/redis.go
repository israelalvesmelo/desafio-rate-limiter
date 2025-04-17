package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"log"

	databasedomain "github.com/israelalvesmelo/desafio-rate-limiter/internal/domain/database"
	"github.com/israelalvesmelo/desafio-rate-limiter/internal/domain/entity"
	"github.com/israelalvesmelo/desafio-rate-limiter/internal/infra/dto"

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
		apiKey.Key,
		jsonReq,
		0,
	).Err(); redisErr != nil {
		log.Println("error inserting rate limit config")
		return redisErr
	}

	return nil
}

func (r *RedisDb) GetRateLimitConfig(ctx context.Context, key string) (*entity.RateLimitConfig, error) {
	jsonReq, err := r.client.Get(ctx, key).Result()
	if err != nil {
		log.Println("error getting rate limit config")
		return nil, err
	}

	var rateLimitConfig entity.RateLimitConfig
	if err := json.Unmarshal([]byte(jsonReq), &rateLimitConfig); err != nil {
		log.Println("error unmarshaling rate limit config")
		return nil, err
	}
	return &rateLimitConfig, nil
}

func (r *RedisDb) UpsertRequest(ctx context.Context, key string, rl *entity.RateLimiter) error {
	req := dto.RequestDB{
		MaxRequests:   rl.MaxRequests,
		TimeWindowSec: rl.TimeWindowSec,
		Requests: func() []int64 {
			reqInt := make([]int64, 0)
			for _, r := range rl.Requests {
				reqInt = append(reqInt, r.Unix())
			}
			return reqInt
		}(),
	}

	jsonReq, marErr := json.Marshal(req)
	if marErr != nil {
		log.Println("error marshaling")
		return marErr
	}

	redisErr := r.client.Set(ctx, createRatePrefix(key), jsonReq, 0).Err()
	if redisErr != nil {
		log.Println("error inserting value")
		return redisErr
	}

	return nil
}

// SaveBlockedDuration Stores the blocked duration amount by key
func (r *RedisDb) SaveBlockedDuration(ctx context.Context, key string, BlockedDuration int64) error {
	if redisErr := r.client.Set(
		ctx,
		createDurationPrefix(key),
		entity.StatusBlocked,
		time.Second*time.Duration(BlockedDuration),
	).Err(); redisErr != nil {
		log.Println("error inserting SaveBlockedDuration")
		return redisErr
	}

	return nil
}

// GetBlockedDuration Obtain the blocked duration by key
func (r *RedisDb) GetBlockedDuration(ctx context.Context, key string) (string, error) {
	val, getErr := r.client.Get(ctx, createDurationPrefix(key)).Result()
	if errors.Is(getErr, redis.Nil) {
		log.Println("INFO: GetBlockedDuration does not exist")
		return "", nil
	}
	if getErr != nil {
		return "", getErr
	}

	return val, nil
}

// GetRequest reads the stored array of request
func (r *RedisDb) GetRequest(ctx context.Context, key string) (*entity.RateLimiter, error) {
	val, getErr := r.client.Get(ctx, createRatePrefix(key)).Result()
	if errors.Is(getErr, redis.Nil) {
		log.Println("INFO: GetRequest key does not exist")
		return &entity.RateLimiter{
			Requests:      make([]time.Time, 0),
			TimeWindowSec: 0,
			MaxRequests:   0,
		}, nil
	}
	if getErr != nil {
		return nil, getErr
	}

	var rateLimiter dto.RequestDB
	if err := json.Unmarshal([]byte(val), &rateLimiter); err != nil {
		log.Println("RateLimiter unmarshal error")
		return &entity.RateLimiter{}, err
	}

	return &entity.RateLimiter{
		Requests: func() []time.Time {
			reqTimeStamp := make([]time.Time, 0)
			for _, rr := range rateLimiter.Requests {
				reqTimeStamp = append(reqTimeStamp, time.Unix(rr, 0))
			}
			return reqTimeStamp
		}(),
		TimeWindowSec: rateLimiter.TimeWindowSec,
		MaxRequests:   rateLimiter.MaxRequests,
	}, nil
}

func createDurationPrefix(key string) string {
	return fmt.Sprintf("%s_%s", entity.PrefixBlockedDurationKey, key)
}

func createRatePrefix(key string) string {
	return fmt.Sprintf("%s_%s", entity.PrefixRateKey, key)
}
