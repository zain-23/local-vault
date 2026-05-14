package main

// lock.go handles "lv lock"
// Removes session from OS keychain immediately
// Forces re-authentication on next command

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/zain-23/local-vault/internal/session"
)

var lockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Lock vault and clear session",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		lvDir := filepath.Join(dir, ".lv")

		// Check if already locked
		if !session.IsUnlocked(lvDir) {
			fmt.Println("🔒 Vault is already locked")
			return nil
		}

		// Delete session for this specific vault only
		// Other project vaults remain unlocked
		if err := session.Delete(lvDir); err != nil {
			return fmt.Errorf("failed to lock vault: %w", err)
		}

		fmt.Println("🔒 Vault locked successfully")
		fmt.Println("   Other project vaults remain unlocked")
		fmt.Println("   Run: lv unlock to unlock again")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(lockCmd)
}
