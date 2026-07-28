package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config is vault-local state stored in .lv/config.json.
// Server URL is machine-global (internal/appstate) — not stored here.
type Config struct {
	WorkspaceID string `json:"workspace_id,omitempty"`
	VaultID     string `json:"vault_id,omitempty"`
	DeviceID    string `json:"device_id,omitempty"`

	// legacySignalingServer is read from old configs for optional heal only.
	// It is never written back by Save.
	legacySignalingServer string `json:"-"`
}

type fileConfig struct {
	WorkspaceID     string `json:"workspace_id"`
	VaultID         string `json:"vault_id"`
	DeviceID        string `json:"device_id"`
	SignalingServer string `json:"signaling_server"`
}

const configFile = "config.json"

func Load(lvDir string) (*Config, error) {
	configPath := filepath.Join(lvDir, configFile)
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}

	var fc fileConfig
	if err := json.Unmarshal(data, &fc); err != nil {
		return nil, err
	}
	return &Config{
		WorkspaceID:           fc.WorkspaceID,
		VaultID:               fc.VaultID,
		DeviceID:              fc.DeviceID,
		legacySignalingServer: fc.SignalingServer,
	}, nil
}

// LegacySignalingServer returns a signaling_server value from an old config
// file, if any. Callers may one-time-heal into appstate; do not use as API base
// when appstate already has a URL.
func (c *Config) LegacySignalingServer() string { return c.legacySignalingServer }

func Save(lvDir string, cfg *Config) error {
	configPath := filepath.Join(lvDir, configFile)
	// saveShape omits signaling_server so it is never written back.
	out := struct {
		WorkspaceID string `json:"workspace_id,omitempty"`
		VaultID     string `json:"vault_id,omitempty"`
		DeviceID    string `json:"device_id,omitempty"`
	}{
		WorkspaceID: cfg.WorkspaceID,
		VaultID:     cfg.VaultID,
		DeviceID:    cfg.DeviceID,
	}
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}
