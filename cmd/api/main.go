package main

import (
	"fmt"

	"github.com/israelalvesmelo/desafio-rate-limiter/cmd/api/config"
	"github.com/israelalvesmelo/desafio-rate-limiter/internal/domain/usecase"
	"github.com/israelalvesmelo/desafio-rate-limiter/internal/infra/database"
	"github.com/israelalvesmelo/desafio-rate-limiter/internal/infra/handler"
	"github.com/israelalvesmelo/desafio-rate-limiter/internal/infra/webserver"
	"github.com/redis/go-redis/v9"
)

func main() {

	// Load config
	var cfg config.Config
	viperCfg := config.NewViper("/config/env")
	viperCfg.ReadViper(&cfg)

	// Create redis client
	redisClient := redis.NewClient(
		&redis.Options{
			Addr: fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
			DB:   cfg.Redis.Db,
		},
	)

	// Create gateway
	storageGateway := database.NewRedisDb(redisClient)

	// Create use case
	createApiKey := usecase.NewRateLimitConfigUseCase(storageGateway)

	server := webserver.NewWebServer(fmt.Sprintf(":%s", cfg.App.Port))
	server.AddHandler("/limiter", handler.NewAPIKeyHandler(createApiKey).CreateAPIKey)
	server.Start()
}
