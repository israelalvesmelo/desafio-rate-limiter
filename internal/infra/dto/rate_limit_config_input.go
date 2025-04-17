package dto

type RateLimitConfigInput struct {
	IsToken       bool  `json:"is_token"`
	MaxRequests   int   `json:"limit" validate:"required"`
	TimeWindow    int64 `json:"time_window" validate:"required"`
	BlockDuration int64 `json:"block_duration" validate:"required"`
}
