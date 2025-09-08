package config

import (
	"os"
	"sync"	
)

type Config struct {
	Endpoint string
}

var cfg Config
var once sync.Once

func loadConfig() {
	//cfg.Endpoint = getEnv("NOVAKEY_ENDPOINT", "https://novakey-api.core-regulus.com")
	cfg.Endpoint = getEnv("NOVAKEY_ENDPOINT", "http://localhost:5000")
}

func Get() *Config {
	once.Do(func() {
		loadConfig()
	})
	return &cfg
}

func getEnv(key string, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}
