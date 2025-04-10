package entity

type RateLimitConfig struct {
	Value         string `json:"value" binding:"required"`          // Ip ou token
	Limit         int    `json:"limit" binding:"required"`          // Máximo de requisições por segundo
	BlockDuration int    `json:"block_duration" binding:"required"` // Duração do bloqueio em segundos
}
