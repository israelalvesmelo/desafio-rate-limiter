package dto

type RateLimitConfigOutput struct {
	Key           string `json:"key"`
	Limit         int    `json:"limit"`
	BlockDuration int    `json:"block_duration"`
}
