package session

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/zalando/go-keyring"
)

const (
	serviceName     = "LocalVault"
	sessionDuration = 12 * time.Hour
)

// Session holds cached session data
type Session struct {
	VaultKey  string    `json:"vault_key"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	VaultID   string    `json:"vault_id"`
}

// Save stores session in OS keychain
// Uses unique key per vault so multiple projects
// can be unlocked simultaneously
func Save(vaultDir string, key []byte) error {
	vid := vaultID(vaultDir)

	session := Session{
		VaultKey:  hex.EncodeToString(key),
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(sessionDuration),
		VaultID:   vid,
	}

	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to serialize session: %w", err)
	}

	// Use vault-specific key name
	// "vault-session-f9a3b2c1" instead of "vault-session"
	// Each project gets its own keychain entry
	keyName := sessionKeyName(vaultDir)

	if err := keyring.Set(serviceName, keyName, string(data)); err != nil {
		// Keychain failed — fall back to temp file
		return saveToFile(vid, data)
	}

	return nil
}

// Load retrieves session from OS keychain
func Load(vaultDir string) ([]byte, error) {
	keyName := sessionKeyName(vaultDir)
	vid := vaultID(vaultDir)

	// Try keychain first
	data, err := keyring.Get(serviceName, keyName)
	if err != nil {
		// Try temp file fallback
		rawData, fileErr := loadFromFile(vid)
		if fileErr != nil {
			return nil, errors.New("vault is locked — run: lv unlock")
		}
		data = string(rawData)
	}

	var session Session
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return nil, errors.New("invalid session — run: lv unlock")
	}

	// Check expired
	if time.Now().After(session.ExpiresAt) {
		_ = Delete(vaultDir)
		return nil, fmt.Errorf(
			"session expired — run: lv unlock\n  (sessions last %s)",
			sessionDuration,
		)
	}

	key, err := hex.DecodeString(session.VaultKey)
	if err != nil {
		return nil, errors.New("corrupted session — run: lv unlock")
	}

	return key, nil
}

// Delete removes session for this specific vault
func Delete(vaultDir string) error {
	keyName := sessionKeyName(vaultDir)
	vid := vaultID(vaultDir)

	// Delete from keychain
	keyring.Delete(serviceName, keyName)

	// Also delete temp file if exists
	deleteFile(vid)

	return nil
}

// IsUnlocked checks if this specific vault is unlocked
func IsUnlocked(vaultDir string) bool {
	_, err := Load(vaultDir)
	return err == nil
}

// TimeRemaining returns how long session has left
func TimeRemaining(vaultDir string) (time.Duration, error) {
	keyName := sessionKeyName(vaultDir)
	vid := vaultID(vaultDir)

	var data string
	var err error

	data, err = keyring.Get(serviceName, keyName)
	if err != nil {
		// Try temp file
		rawData, fileErr := loadFromFile(vid)
		if fileErr != nil {
			return 0, errors.New("vault is locked")
		}
		data = string(rawData)
	}

	var session Session
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return 0, errors.New("invalid session")
	}

	remaining := time.Until(session.ExpiresAt)
	if remaining < 0 {
		return 0, errors.New("session expired")
	}

	return remaining, nil
}

// sessionKeyName returns unique keychain key for this vault
// "vault-session-f9a3b2c1" — different per project
func sessionKeyName(vaultDir string) string {
	return fmt.Sprintf("vault-session-%s", vaultID(vaultDir))
}

// vaultID creates unique identifier for a vault directory
// Always uses absolute path for consistency
func vaultID(vaultDir string) string {
	absPath, err := filepath.Abs(vaultDir)
	if err != nil {
		absPath = vaultDir
	}
	hash := sha256.Sum256([]byte(absPath))
	return hex.EncodeToString(hash[:])[:16]
}

// ===== TEMP FILE FALLBACK =====
// Used when keychain is unavailable

func sessionFilePath(vid string) string {
	uid := os.Getuid()
	return filepath.Join(
		os.TempDir(),
		fmt.Sprintf(".lv-session-%d-%s", uid, vid),
	)
}

func saveToFile(vid string, data []byte) error {
	path := sessionFilePath(vid)
	return os.WriteFile(path, data, 0600)
}

func loadFromFile(vid string) ([]byte, error) {
	path := sessionFilePath(vid)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("no session file")
		}
		return nil, err
	}
	return data, nil
}

func deleteFile(vid string) {
	path := sessionFilePath(vid)
	os.Remove(path)
}
