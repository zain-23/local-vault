package config

// config.go stores LocalVault configuration
// Things like: which signaling server to use

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds LocalVault settings
type Config struct {
	// SignalingServer is the URL of the signaling server
	// Default: our hosted server
	// Teams can self-host and change this
	SignalingServer string `json:"signaling_server"`

	// DeviceID cached here for quick access
	DeviceID string `json:"device_id"`
}

// Default server URL
// Change this to your deployed server URL later
const DefaultServer = "http://localhost:8080"

const configFile = "config.json"

// Load reads config from .lv/config.json
// Returns default config if file not found
func Load(lvDir string) (*Config, error) {
	configPath := filepath.Join(lvDir, configFile)

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default config if no config file yet
			return &Config{
				SignalingServer: DefaultServer,
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

// Save writes config to .lv/config.json
func Save(lvDir string, cfg *Config) error {
	configPath := filepath.Join(lvDir, configFile)

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
