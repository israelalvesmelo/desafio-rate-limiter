package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"

	"github.com/israelalvesmelo/desafio-rate-limiter/internal/domain/database"
	"github.com/israelalvesmelo/desafio-rate-limiter/internal/domain/entity"
	"github.com/israelalvesmelo/desafio-rate-limiter/internal/infra/dto"
)

type CreateRateLimitConfigInput struct {
	Ip            string `json:"ip"`             // Ip ou token
	IsToken       bool   `json:"is_token"`       // Se é um token ou ip
	Limit         int    `json:"limit"`          // Máximo de requisições por segundo
	BlockDuration int    `json:"block_duration"` // Duração do bloqueio em segundos
}

type CreateRateLimitConfigUseCase interface {
	Execute(ctx context.Context, input CreateRateLimitConfigInput) (*dto.RateLimitConfigOutput, error)
}

type CreateRateLimitConfigUseCaseImpl struct {
	StorageGateway database.StorageDb
}

func NewRateLimitConfigUseCase(storageGateway database.StorageDb) CreateRateLimitConfigUseCase {
	return &CreateRateLimitConfigUseCaseImpl{
		StorageGateway: storageGateway,
	}
}

func (c *CreateRateLimitConfigUseCaseImpl) Execute(ctx context.Context, input CreateRateLimitConfigInput) (*dto.RateLimitConfigOutput, error) {
	rtConf := &entity.RateLimitConfig{
		Limit:         input.Limit,
		BlockDuration: input.BlockDuration,
	}

	var value string
	var err error

	if input.IsToken {
		value, err = c.generateTokenValue()
		if err != nil {
			log.Println("error generating token value")
			return nil, err
		}
	} else {
		value = input.Ip
	}

	rtConf.Value = value
	err = c.StorageGateway.SaveRateLimitConfig(ctx, rtConf)
	if err != nil {
		return nil, err
	}

	return &dto.RateLimitConfigOutput{
		Key: rtConf.Value,
	}, nil
}

func (c *CreateRateLimitConfigUseCaseImpl) generateTokenValue() (string, error) {
	bytes, err := c.generateRandomBytes(32)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (c *CreateRateLimitConfigUseCaseImpl) generateRandomBytes(length int) ([]byte, error) {
	byteSlice := make([]byte, length)
	_, err := rand.Read(byteSlice)
	if err != nil {
		return nil, err
	}

	return byteSlice, nil
}
