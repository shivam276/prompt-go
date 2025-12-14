package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Anthropic AnthropicConfig `yaml:"anthropic"`
}

type AnthropicConfig struct {
	APIKey string `yaml:"api_key"`
	Model  string `yaml:"model"` // "claude-3-5-haiku-20241022" or "claude-3-5-sonnet-20241022"
}

// Load loads configuration from ~/.promptgo/config.yaml or environment variables
func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(home, ".promptgo", "config.yaml")

	data, err := os.ReadFile(configPath)
	if err != nil {
		// Return default config if file doesn't exist
		if os.IsNotExist(err) {
			return &Config{
				Anthropic: AnthropicConfig{
					Model: "claude-3-5-haiku-20241022",
				},
			}, nil
		}
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Check for API key in environment variable as fallback
	if cfg.Anthropic.APIKey == "" {
		cfg.Anthropic.APIKey = os.Getenv("ANTHROPIC_API_KEY")
	}

	// Set default model if not specified
	if cfg.Anthropic.Model == "" {
		cfg.Anthropic.Model = "claude-3-5-haiku-20241022"
	}

	return &cfg, nil
}
