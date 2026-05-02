package cmd

// unlock.go handles "lv unlock"
// Asks passphrase once and caches in OS keychain
// All other commands read from cache — zero friction
//
// Like: ssh-add (adds SSH key to agent so no password needed)
// Like: op signin (1Password CLI unlock)

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/zain-23/local-vault/internal/session"
	"github.com/zain-23/local-vault/internal/vault"
)

var unlockCmd = &cobra.Command{
	Use:   "unlock",
	Short: "Unlock vault for this session (12 hours)",
	Long: `Unlocks the vault by asking your passphrase once.
The derived key is cached in your OS keychain.
All other commands work without prompting until session expires.

Session lasts 12 hours or until you run: lv lock`,

	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		lvDir := filepath.Join(dir, ".lv")

		// Check if already unlocked
		if session.IsUnlocked(lvDir) {
			remaining, _ := session.TimeRemaining(lvDir)
			hours := int(remaining.Hours())
			minutes := int(remaining.Minutes()) % 60

			fmt.Printf("✅ Already unlocked (%dh %dm remaining)\n",
				hours, minutes)
			fmt.Println("   Run: lv lock to lock now")
			return nil
		}

		// Ask passphrase
		fmt.Println("🔐 Unlocking LocalVault...")
		passphrase, err := promptPassphrase()
		if err != nil {
			return err
		}

		// Load vault to verify passphrase is correct
		// This also derives the key
		v, err := vault.Load(dir, passphrase)
		if err != nil {
			return err
		}

		// Get derived key
		key := v.GetKey()

		// Save key to OS keychain
		if err := session.Save(lvDir, key); err != nil {
			return fmt.Errorf("failed to save session: %w", err)
		}

		// Calculate expiry time
		expiresAt := time.Now().Add(12 * time.Hour)

		fmt.Println()
		fmt.Println("✅ Vault unlocked successfully")
		fmt.Printf("⏱️  Session valid until: %s\n",
			expiresAt.Format("15:04:05 (Jan 02)"))
		fmt.Println()
		fmt.Println("You can now run any lv command without entering passphrase.")
		fmt.Println("Run: lv lock   to lock immediately")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(unlockCmd)
}
