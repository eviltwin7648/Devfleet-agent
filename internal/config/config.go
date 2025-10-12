package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	APIKey string `json:"api_key`
}

func ConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "./devfleet", "config.json")
}

func SaveKey(key string) error {
	cfg := Config{APIKey: key}
	data, err := json.MarshalIndent(cfg, "", " ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(ConfigPath()), 0755); err != nil {
		return err
	}
	return os.WriteFile(ConfigPath(), data, 0600)
}

func LoadKey() (*Config, error) {
	data, err := os.ReadFile(ConfigPath())

	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
