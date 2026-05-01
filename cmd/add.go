package cmd

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/zain-23/local-vault/internal/vault"
	"golang.org/x/term"
)

// env flag value — shared across commands
var envFlag string

var addCmd = &cobra.Command{
	Use:   "add KEY=VALUE",
	Short: "Add or update a secret",
	Example: `  lv add DATABASE_URL=postgres://localhost/mydb
  lv add API_KEY=sk-xxx --env production
  lv add SECRET_KEY=abc123 --env staging`,
	Args: cobra.ExactArgs(1), // requires exactly 1 argument
	RunE: func(cmd *cobra.Command, args []string) error {
		// Parse KEY=VALUE argument
		parts := strings.SplitN(args[0], "=", 2) // split on first = only
		if len(parts) != 2 {
			return fmt.Errorf("invalid format. use: lv add KEY=VALUE")
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "" {
			return fmt.Errorf("key cannot be empty")
		}

		// Load vault (needs passphrase)
		passphrase, err := promptPassphrase()
		if err != nil {
			return err
		}

		dir, _ := os.Getwd()
		v, err := vault.Load(dir, passphrase)
		if err != nil {
			return err
		}

		// Add the secret
		if err := v.Add(key, value, envFlag); err != nil {
			return err
		}

		env := envFlag
		if env == "" {
			env = "all environments"
		}

		fmt.Printf("✅ Added %s (%s)\n", key, env)
		return nil
	},
}

func init() {
	// Add --env flag to add command
	// Like: program.option('--env <env>', 'environment') in commander.js
	addCmd.Flags().StringVarP(&envFlag, "env", "e", "", "environment (development/staging/production)")
	rootCmd.AddCommand(addCmd)
}

// promptPassphrase asks user for passphrase without showing it on screen
// Reused across multiple commands
func promptPassphrase() (string, error) {
	fmt.Print("Passphrase: ")
	passphraseBytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", fmt.Errorf("failed to read passphrase: %w", err)
	}
	return string(passphraseBytes), nil
}
