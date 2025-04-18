package config

// DataBase configuration related to this specific database
type DataBase struct {
	Name string
	Db   int
	Host string
	Port string
}

// App General configuration related to the Application
type App struct {
	Host string
	Port string
}

// RateLimiter properties for configuration
type RateLimiter struct {
	ByIP     LimitValues
	ByAPIKey LimitValues
}

type LimitValues struct {
	MaxRequests     int
	TimeWindow      int64
	BlockedDuration int64
}

// Config Final Struct Configuration
type Config struct {
	DataBase    DataBase
	App         App
	RateLimiter RateLimiter
}
