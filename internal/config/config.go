package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

var DefaultServer = "http://localhost:8080"

type Config struct {
	SignalingServer string `json:"signaling_server"`
	DeviceID        string `json:"device_id"`
	VaultID         string `json:"vault_id"` // ← add this
}

const configFile = "config.json"

func Load(lvDir string) (*Config, error) {
	serverURL := os.Getenv("SERVER_URL")

	configPath := filepath.Join(lvDir, configFile)
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{
				SignalingServer: serverURL,
			}, nil
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func Save(lvDir string, cfg *Config) error {
	configPath := filepath.Join(lvDir, configFile)
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}
