package main

import (
	"github.com/israelalvesmelo/desafio-rate-limiter/internal/domain/usecase"
	"github.com/israelalvesmelo/desafio-rate-limiter/internal/infra/database"
	"github.com/israelalvesmelo/desafio-rate-limiter/internal/infra/handler"
	"github.com/israelalvesmelo/desafio-rate-limiter/internal/infra/webserver"
)

func main() {

	// Load config
	/*
			var cfg config.Config
		viperCfg := config.NewViper("env")
		viperCfg.ReadViper(&cfg)
	*/

	// Create redis client
	/*	redisClient := redis.NewClient(
			&redis.Options{
				Addr: "",
				DB:   nil,
			},
		),
	*/
	// Create gateway
	storageGateway := database.NewRedisDb(nil)

	// Create use case
	createApiKey := usecase.NewRateLimitConfigUseCase(storageGateway)

	server := webserver.NewWebServer(":8080")
	server.AddHandler("/api/v1/api-key", handler.NewAPIKeyHandler(createApiKey).CreateAPIKey)
	server.Start()
}

//TODO: REMOVER DAQUI
