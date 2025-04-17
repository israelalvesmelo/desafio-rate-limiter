package dto

type RateLimitConfigOutput struct {
	Key             string `json:"key"`
	MaxRequests     int    `json:"max_requests"`
	TimeWindow      int64  `json:"time_window"`
	BlockedDuration int64  `json:"block_duration"`
}
