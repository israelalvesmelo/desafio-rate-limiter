package dto

type RateLimitConfigInput struct {
	IsToken       bool `json:"is_token"`       // Se é um token ou ip
	Limit         int  `json:"limit"`          // Máximo de requisições por segundo
	BlockDuration int  `json:"block_duration"` // Duração do bloqueio em segundos
}
