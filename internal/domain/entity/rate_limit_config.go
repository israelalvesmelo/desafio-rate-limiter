package entity

const (
	PrefixRateKey            = "rate:key"
	PrefixBlockedDurationKey = "blocked:key"
	StatusBlocked            = "Blocked"
	APIKeyName               = "API_KEY"
	IPName                   = "IP"
)

type RateLimitConfig struct {
	Key string `json:"key" binding:"required"`
	LimitValues
}

type LimitValues struct {
	MaxRequests     int   `json:"max_requests" binding:"required"`
	TimeWindow      int64 `json:"time_window" binding:"required"`
	BlockedDuration int64 `json:"block_duration" binding:"required"`
}
