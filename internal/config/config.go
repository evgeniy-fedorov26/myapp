package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config - структура конфигурации.
type Config struct {
	Rss           []string `json:"rss"`
	RequestPeriod int      `json:"request_period"`
	DatabaseURL   string   `json:"database_url"`
}

// LoadConfig загружает конфигурацию из файла.
func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %w", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("error decoding config file: %w", err)
	}

	// Валидация конфигурации
	if len(config.Rss) == 0 {
		return nil, fmt.Errorf("config error: no RSS feeds specified")
	}
	if config.RequestPeriod <= 0 {
		return nil, fmt.Errorf("config error: request_period must be a positive integer")
	}
	if config.DatabaseURL == "" {
		return nil, fmt.Errorf("config error: database_url is required")
	}

	return &config, nil
}
