package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/zain-23/local-vault/internal/session"
	"github.com/zain-23/local-vault/internal/vault"
	"golang.org/x/term"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new vault in the current directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔐 Initializing LocalVault...")

		// Ask passphrase
		fmt.Print("Enter passphrase: ")
		passphraseBytes, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			return fmt.Errorf("failed to read passphrase: %w", err)
		}

		// Confirm passphrase
		fmt.Print("Confirm passphrase: ")
		confirmBytes, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			return fmt.Errorf("failed to read passphrase: %w", err)
		}

		if string(passphraseBytes) != string(confirmBytes) {
			return fmt.Errorf("passphrases do not match")
		}

		passphrase := string(passphraseBytes)

		if len(passphrase) < 8 {
			return fmt.Errorf("passphrase must be at least 8 characters")
		}

		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		// Initialize vault on disk
		if err := vault.Init(dir, passphrase); err != nil {
			return err
		}

		// Auto unlock after init
		// User just proved they know the passphrase
		// No need to run lv unlock separately right after lv init
		v, err := vault.Load(dir, passphrase)
		if err == nil {
			lvDir := filepath.Join(dir, ".lv")
			if err := session.Save(lvDir, v.GetKey()); err == nil {
				fmt.Println("🔓 Auto-unlocked for 12 hours")
			}
		}

		fmt.Println("✅ Vault initialized successfully")
		fmt.Println("📁 Created .lv/ directory")
		fmt.Println("🔒 Vault encrypted with your passphrase")
		fmt.Println("📝 Updated .gitignore")
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Println("  lv add DATABASE_URL=postgres://...")
		fmt.Println("  lv add API_KEY=sk-xxx")
		fmt.Println("  lv inject -- npm run dev")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
