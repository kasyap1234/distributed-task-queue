package config

import (
	"log"

	"github.com/spf13/viper"
)
type Config struct {
	RedisAddr string 
	MaxRetries int 
	WorkerCount int 
	APIAddr string 
	WorkerTimeout int 

}
func LoadConfig() *Config{
	viper.SetDefault("RedisAddr", "localhost:6379")
	viper.SetDefault("MaxRetries", 5)
	viper.SetDefault("WorkerCount", 10)
	viper.SetDefault("APIAddr", "localhost:8080")
	viper.SetDefault("WorkerTimeout", 10)
	viper.AutomaticEnv()
	config :=&Config{
		RedisAddr: viper.GetString("RedisAddr"),
		MaxRetries: viper.GetInt("MaxRetries"),
		WorkerCount: viper.GetInt("WorkerCount"),
		APIAddr: viper.GetString("APIAddr"),
		WorkerTimeout: viper.GetInt("WorkerTimeout"),
	}
	return config; 
	
}