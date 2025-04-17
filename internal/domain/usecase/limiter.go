package usecase

import (
	"context"
	"log"

	"github.com/israelalvesmelo/desafio-rate-limiter/internal/domain/database"
	"github.com/israelalvesmelo/desafio-rate-limiter/internal/domain/entity"
	"github.com/israelalvesmelo/desafio-rate-limiter/internal/infra/dto"
)

type Limiter struct {
	storage database.StorageDb
}

func NewLimiter(storage database.StorageDb) *Limiter {
	return &Limiter{
		storage: storage,
	}
}

func (l *Limiter) Execute(ctx context.Context, input dto.RequestSave, rtConfig entity.RateLimitConfig) (*dto.RequestResult, error) {
	status, blockedErr := l.storage.GetBlockedDuration(ctx, rtConfig.Key)
	if blockedErr != nil {
		return nil, blockedErr
	}
	if status == entity.StatusBlocked {
		log.Println("ip/key is blocked due to exceeding the maximum number of requests")
		return nil, entity.ErrIPExceededAmountRequest
	}

	getRequest, getReqErr := l.storage.GetRequest(ctx, rtConfig.Key)
	if getReqErr != nil {
		log.Printf("Error getting ip/key requests: %s \n", getReqErr.Error())
		return nil, getReqErr
	}

	getRequest.TimeWindowSec = rtConfig.TimeWindow
	getRequest.MaxRequests = rtConfig.MaxRequests
	if valErr := getRequest.Validate(); valErr != nil {
		log.Printf("Error validation in rate limiter: %s \n", valErr.Error())
		return nil, valErr
	}

	getRequest.AddRequests(input.TimeAdd)
	isAllowed := getRequest.Allow(input.TimeAdd)
	if upsertErr := l.storage.UpsertRequest(ctx, rtConfig.Key, getRequest); upsertErr != nil {
		log.Printf("Error updating/inserting rate limit: %s \n", upsertErr.Error())
		return nil, upsertErr
	}

	if !isAllowed {
		if saveErr := l.storage.SaveBlockedDuration(
			ctx,
			rtConfig.Key,
			rtConfig.BlockedDuration,
		); saveErr != nil {
			return nil, saveErr
		}
	}

	return &dto.RequestResult{
		Allow: isAllowed,
	}, nil
}
