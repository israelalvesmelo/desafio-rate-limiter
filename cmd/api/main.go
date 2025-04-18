package main

import (
	"fmt"

	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/israelalvesmelo/desafio-rate-limiter/cmd/api/config"
	"github.com/israelalvesmelo/desafio-rate-limiter/internal/domain/usecase"
	"github.com/israelalvesmelo/desafio-rate-limiter/internal/infra/database"
	"github.com/israelalvesmelo/desafio-rate-limiter/internal/infra/handler"
	"github.com/israelalvesmelo/desafio-rate-limiter/internal/infra/middleware"

	"github.com/israelalvesmelo/desafio-rate-limiter/internal/infra/webserver"
)

func main() {

	// Load config
	var cfg config.Config
	viperCfg := config.NewViper("env")
	viperCfg.ReadViper(&cfg)

	// Create gateway
	storageGateway, _ := database.GetDatabase(cfg.DataBase)

	// Create use case
	createApiKey := usecase.NewRateLimitConfigUseCase(storageGateway)

	// Create handler
	helloWorldHandler := handler.NewHelloWorldHandler()
	rateLimitConfigHandler := handler.NewRateLimitConfigHandler(createApiKey)

	// Create middleware
	limiter := middleware.NewRateLimitMiddleware(storageGateway, cfg.RateLimiter)

	// Create webserver
	server := webserver.NewWebServer(fmt.Sprintf(":%s", cfg.App.Port))
	server.AddMiddleware(limiter.Handler)
	server.AddMiddleware(chimiddleware.Logger)

	server.AddHandler("/rate-limiter", rateLimitConfigHandler.Create)
	server.AddHandler("/hello-world", helloWorldHandler.HelloWorld)

	server.Start()
}
