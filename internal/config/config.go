package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type RateLimitConfig struct {
	MaxTokens  int           `json:"maxTokens"`
	RefillRate time.Duration `json:"refillRate"`
}

type LoggingConfig struct {
	Frequency time.Duration `json:"frequency"`
}

type Config struct {
	RateLimit RateLimitConfig `json:"rateLimit"`
	Logging   LoggingConfig   `json:"logging"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &Config{}
	if err := decoder.Decode(config); err != nil {
		return nil, fmt.Errorf("could not decode config JSON: %w", err)
	}

	// Convert refill rate to duration
	config.RateLimit.RefillRate, err = time.ParseDuration(config.RateLimit.RefillRate.String())
	if err != nil {
		return nil, fmt.Errorf("invalid refill rate duration: %w", err)
	}

	config.Logging.Frequency, err = time.ParseDuration(config.Logging.Frequency.String())
	if err != nil {
		return nil, fmt.Errorf("invalid logging frequency duration: %w", err)
	}

	return config, nil
}
