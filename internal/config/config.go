package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port       string
	DBPath     string
	MemLimit   string
}

func Load() *Config {
	return &Config{
		Port:     getEnv("PORT", "8822"),
		DBPath:   getEnv("DB_PATH", "./data/go-minitrackr.db"),
		MemLimit: getEnv("GOMEMLIMIT", "25MiB"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func (c *Config) PortInt() int {
	port, _ := strconv.Atoi(c.Port)
	return port
}
