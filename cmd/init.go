package cmd

import (
	"fmt"
	"os"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/zain-23/local-vault/internal/vault"
	"golang.org/x/term" // reads password input without showing on screen
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new vault in the current directory",
	// RunE is like Run but returns an error
	// Cobra handles printing the error for us
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔐 Initializing LocalVault...")

		// Prompt for passphrase (hidden input — like a password field)
		fmt.Print("Enter passphrase: ")
		passphraseBytes, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println() // newline after hidden input
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

		// Check they match
		if string(passphraseBytes) != string(confirmBytes) {
			return fmt.Errorf("passphrases do not match")
		}

		passphrase := string(passphraseBytes)

		if len(passphrase) < 8 {
			return fmt.Errorf("passphrase must be at least 8 characters")
		}

		// Get current directory
		// Like: process.cwd() in Node.js
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		// Initialize vault
		if err := vault.Init(dir, passphrase); err != nil {
			return err
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

// init() in Go runs automatically when package loads
// We use it to register commands with the root command
// Like: app.use('/init', initRouter) in Express
func init() {
	rootCmd.AddCommand(initCmd)
}
