package dto

import "time"

type RequestSave struct {
	TimeAdd time.Time
}

type RequestResult struct {
	Allow bool
}

type RequestDB struct {
	MaxRequests   int     `json:"max_requests"`
	TimeWindowSec int64   `json:"time_window_sec"`
	Requests      []int64 `json:"requests"`
}
