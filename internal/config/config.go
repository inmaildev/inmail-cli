package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	APIKey  string `json:"api_key"`
	BaseURL string `json:"base_url"`
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".inmail", "config.json"), nil
}

func Load() (*Config, error) {
	cfg := &Config{
		BaseURL: "https://inmail.dev/v1",
	}

	if key := os.Getenv("INMAIL_API_KEY"); key != "" {
		cfg.APIKey = key
	}
	if url := os.Getenv("INMAIL_BASE_URL"); url != "" {
		cfg.BaseURL = url
	}

	path, err := configPath()
	if err != nil {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, nil
	}

	var fileCfg Config
	if err := json.Unmarshal(data, &fileCfg); err != nil {
		return cfg, nil
	}

	if cfg.APIKey == "" && fileCfg.APIKey != "" {
		cfg.APIKey = fileCfg.APIKey
	}
	if fileCfg.BaseURL != "" && os.Getenv("INMAIL_BASE_URL") == "" {
		cfg.BaseURL = fileCfg.BaseURL
	}

	return cfg, nil
}

func Save(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}
