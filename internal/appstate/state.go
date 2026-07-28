package appstate

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// defaultServerURL can be overridden at build time via
// -ldflags "-X .../internal/appstate.defaultServerURL=https://...".
var defaultServerURL = "http://localhost:8080"

const stateFile = "state.json"

// State is machine-global CLI state, stored outside any .lv directory because
// login happens before a vault exists.
type State struct {
	DeviceFingerprint string `json:"device_fingerprint"`
	DeviceName        string `json:"device_name"`
	ServerURL         string `json:"server_url"`
}

// Dir returns the CLI config directory (e.g. ~/.config/local-vault).
func Dir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "local-vault"), nil
}

// Load reads state.json, healing/seeding any missing fields, and persists changes.
func Load() (*State, error) {
	dir, err := Dir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dir, stateFile)

	var s State
	data, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		// fresh install — fall through with a zero State to seed below
	} else if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}

	changed := false
	if s.DeviceFingerprint == "" {
		s.DeviceFingerprint = uuid.New().String()
		changed = true
	}
	if s.DeviceName == "" {
		s.DeviceName = hostname()
		changed = true
	}
	if resolved := resolveServerURL(); s.ServerURL != resolved && (s.ServerURL == "" || os.Getenv("SERVER_URL") != "") {
		s.ServerURL = resolved
		changed = true
	}

	if changed {
		if err := save(path, &s); err != nil {
			return nil, err
		}
	}
	return &s, nil
}

func resolveServerURL() string {
	if v := os.Getenv("SERVER_URL"); v != "" {
		return v
	}
	return defaultServerURL
}

func hostname() string {
	h, err := os.Hostname()
	if err != nil || h == "" {
		return "unknown-device"
	}
	return h
}

func save(path string, s *State) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}
