package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

// Config aggregates configuration values for the backend services.
type Config struct {
	Environment      string
	HTTPPort         int
	QueueBufferSize  int
	StorageDSN       string
	JWTSigningSecret string
}

var (
	cfg  *Config
	once sync.Once
)

// FromEnv reads configuration values from environment variables.
func FromEnv() (*Config, error) {
	var err error
	once.Do(func() {
		cfg = &Config{}
		if err = populate(cfg); err != nil {
			cfg = nil
		}
	})
	if cfg == nil {
		return nil, err
	}
	return cfg, nil
}

func populate(cfg *Config) error {
	cfg.Environment = getString("APP_ENV", "development")
	cfg.StorageDSN = getString("STORAGE_DSN", "memory://default")
	cfg.JWTSigningSecret = getString("JWT_SIGNING_SECRET", "dev-secret")

	port, err := strconv.Atoi(getString("HTTP_PORT", "8080"))
	if err != nil {
		return fmt.Errorf("invalid HTTP_PORT: %w", err)
	}
	cfg.HTTPPort = port

	bufferSize, err := strconv.Atoi(getString("QUEUE_BUFFER_SIZE", "100"))
	if err != nil {
		return fmt.Errorf("invalid QUEUE_BUFFER_SIZE: %w", err)
	}
	cfg.QueueBufferSize = bufferSize

	return nil
}

func getString(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
