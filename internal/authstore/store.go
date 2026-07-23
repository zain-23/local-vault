package authstore

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/zalando/go-keyring"

	"github.com/zain-23/local-vault/internal/appstate"
)

const (
	serviceName = "LocalVault"
	keyName     = "auth"
	authFile    = "auth.json"
)

// ErrNoTokens means the user is not logged in.
var ErrNoTokens = errors.New("not logged in")

// Tokens is the CLI's session — distinct from the per-vault crypto session.
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	DeviceID     string `json:"device_id"`
}

// Save writes tokens to the OS keychain, falling back to a 0600 file.
func Save(t *Tokens) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	if err := keyring.Set(serviceName, keyName, string(data)); err != nil {
		return saveToFile(data)
	}
	return nil
}

// Load reads tokens from the keychain, falling back to the file.
func Load() (*Tokens, error) {
	data, err := keyring.Get(serviceName, keyName)
	if err != nil {
		raw, ferr := loadFromFile()
		if ferr != nil {
			return nil, ErrNoTokens
		}
		data = string(raw)
	}
	var t Tokens
	if err := json.Unmarshal([]byte(data), &t); err != nil {
		return nil, ErrNoTokens
	}
	if t.AccessToken == "" {
		return nil, ErrNoTokens
	}
	return &t, nil
}

// Clear removes tokens from both the keychain and the file fallback.
func Clear() error {
	_ = keyring.Delete(serviceName, keyName)
	if path, err := filePath(); err == nil {
		_ = os.Remove(path)
	}
	return nil
}

func filePath() (string, error) {
	dir, err := appstate.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, authFile), nil
}

func saveToFile(data []byte) error {
	path, err := filePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func loadFromFile() ([]byte, error) {
	path, err := filePath()
	if err != nil {
		return nil, err
	}
	return os.ReadFile(path)
}
