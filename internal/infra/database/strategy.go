package database

import (
	"fmt"

	"github.com/israelalvesmelo/desafio-rate-limiter/cmd/api/config"
	databasedomain "github.com/israelalvesmelo/desafio-rate-limiter/internal/domain/database"
	"github.com/redis/go-redis/v9"
)

// GetDatabase returns the database strategy based on the configuration
func GetDatabase(cfg config.DataBase) (databasedomain.StorageDb, error) {
	switch cfg.Name {
	default:
		redisClient := redis.NewClient(
			&redis.Options{
				Addr: fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
				DB:   cfg.Db,
			},
		)
		return newRedisDb(redisClient), nil
	}
}
