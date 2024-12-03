package config

import (
	"os"
	"time"
)

type Config struct {
	Database DataBaseConfig
	GRPC     GRPCConfig
}

type DataBaseConfig struct {
	URL           string
	MaxRetries    int
	RetryInterval time.Duration
}

type GRPCConfig struct {
	Port string
}

func Load() *Config {
	return &Config{
		Database: DataBaseConfig{
			URL:           os.Getenv("DATABASE_URL"),
			MaxRetries:    5,
			RetryInterval: 5 * time.Second,
		},
		GRPC: GRPCConfig{
			Port: genEnvOrDefault("GRPC_PORT", "50051"),
		},
	}
}

func genEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
