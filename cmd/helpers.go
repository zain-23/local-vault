package cmd

import (
	"fmt"
	"path/filepath"
	"syscall"

	"github.com/zain-23/local-vault/internal/session"
	"github.com/zain-23/local-vault/internal/vault"
	"golang.org/x/term"
)

var envFlag string

func promptPassphrase() (string, error) {
	fmt.Print("Passphrase: ")
	passphraseBytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", fmt.Errorf("failed to read passphrase: %w", err)
	}
	if len(passphraseBytes) == 0 {
		return "", fmt.Errorf("passphrase cannot be empty")
	}
	return string(passphraseBytes), nil
}

// loadVault loads vault using session key if available
// No passphrase needed after lv unlock
func loadVault(dir string) (*vault.Vault, error) {
	lvDir := filepath.Join(dir, ".lv")

	// Try session cache first
	key, err := session.Load(lvDir)
	if err == nil {
		return vault.LoadWithKey(dir, key)
	}

	// Session not available — tell user to unlock
	return nil, fmt.Errorf(
		"vault is locked\n\n  Run: lv unlock\n  (unlocks for 12 hours)",
	)
}
