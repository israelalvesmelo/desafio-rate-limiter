package entity

import "errors"

var (
	ErrExceededAmountRequest  = errors.New("you have reached the maximum number of Requests or actions allowed within a certain time frame")
	ErrBlockedTimeDuration    = errors.New("blocked time duration should be greater than zero")
	ErrRateLimiterTimeWindow  = errors.New("rate limiter time window duration should be greater than zero")
	ErrRateLimiterMaxRequests = errors.New("rate limiter maximum requests should be greater than zero")
)
