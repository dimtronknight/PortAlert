package stocks

import (
	"fmt"
	"os"
)

type Config struct {
	APIKey string
	IsLive bool
}

func LoadConfig() (*Config, error) {
	apiKey := os.Getenv("TRADING212_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("TRADING212_API_KEY environment variable is required")
	}

	isLive := os.Getenv("TRADING212_IS_LIVE") == "true"

	return &Config{
		APIKey: apiKey,
		IsLive: isLive,
	}, nil
}

func NewClientFromConfig() (*Client, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	return NewClient(config.APIKey, config.IsLive), nil
}