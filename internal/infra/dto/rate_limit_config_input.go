package dto

type RateLimitConfigInput struct {
	IsToken       bool `json:"is_token"`                           // Se é um token ou ip
	Limit         int  `json:"limit" validate:"required"`          // Máximo de requisições por segundo
	BlockDuration int  `json:"block_duration" validate:"required"` // Duração do bloqueio em segundos
}
