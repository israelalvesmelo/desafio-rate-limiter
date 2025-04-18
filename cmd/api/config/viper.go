package config

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const fileExtension = "json"

type Viper struct {
	fileName string
}

func NewViper(fileName string) *Viper {
	return &Viper{
		fileName: fileName,
	}
}

func (v *Viper) ReadViper(config *Config) {
	viper.SetConfigName(v.fileName)
	viper.SetConfigType(fileExtension)
	viper.AddConfigPath(".")
	viper.WatchConfig()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("fatal error config file: default \n", err)
		os.Exit(1)
	}

	v.readConfig(config)
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		v.readConfig(config)
	})
	v.logConfig(config)
}

func (v *Viper) readConfig(c *Config) {
	c.DataBase.Name = viper.GetString("database.name")
	c.DataBase.Db = viper.GetInt("database.db")
	c.DataBase.Host = viper.GetString("database.host")
	c.DataBase.Port = viper.GetString("database.port")

	c.App.Host = viper.GetString("app.host")
	c.App.Port = viper.GetString("app.port")

	c.RateLimiter.ByIP.BlockedDuration = viper.GetInt64("rate_limiter.by_ip.blocked_duration")
	c.RateLimiter.ByIP.TimeWindow = viper.GetInt64("rate_limiter.by_ip.time_window")
	c.RateLimiter.ByIP.MaxRequests = viper.GetInt("rate_limiter.by_ip.max_requests")

	c.RateLimiter.ByAPIKey.BlockedDuration = viper.GetInt64("rate_limiter.by_api_key.blocked_duration")
	c.RateLimiter.ByAPIKey.TimeWindow = viper.GetInt64("rate_limiter.by_api_key.time_window")
	c.RateLimiter.ByAPIKey.MaxRequests = viper.GetInt("rate_limiter.by_api_key.max_requests")
}

func (v *Viper) logConfig(c *Config) {
	fmt.Println("=== Configuration Loaded ===")

	fmt.Println("\n[DataBase]")
	fmt.Printf("Name: %s\n", c.DataBase.Name)
	fmt.Printf("DB: %d\n", c.DataBase.Db)
	fmt.Printf("Host: %s\n", c.DataBase.Host)
	fmt.Printf("Port: %s\n", c.DataBase.Port)

	fmt.Println("\n[App]")
	fmt.Printf("Host: %s\n", c.App.Host)
	fmt.Printf("Port: %s\n", c.App.Port)

	fmt.Println("\n[Rate Limiter - By IP]")
	fmt.Printf("Blocked Duration: %d seconds\n", c.RateLimiter.ByIP.BlockedDuration)
	fmt.Printf("Time Window: %d seconds\n", c.RateLimiter.ByIP.TimeWindow)
	fmt.Printf("Max Requests: %d\n", c.RateLimiter.ByIP.MaxRequests)

	fmt.Println("\n[Rate Limiter - By API Key]")
	fmt.Printf("Blocked Duration: %d seconds\n", c.RateLimiter.ByAPIKey.BlockedDuration)
	fmt.Printf("Time Window: %d seconds\n", c.RateLimiter.ByAPIKey.TimeWindow)
	fmt.Printf("Max Requests: %d\n", c.RateLimiter.ByAPIKey.MaxRequests)

	fmt.Println("\n=== End of Configuration ===")
}
