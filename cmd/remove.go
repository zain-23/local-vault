package cmd

// remove.go handles "lv remove KEY"
// Deletes a secret from the vault permanently

import (
	"bufio" // reads user input line by line
	"fmt"
	"os"
	"strings" // string manipulation (like .trim() in JS)

	"github.com/spf13/cobra"
	"github.com/zain-23/local-vault/internal/vault"
)

var removeCmd = &cobra.Command{
	Use:   "remove KEY",
	Short: "Remove a secret from the vault",
	Example: `  lv remove DATABASE_URL
  lv remove STRIPE_KEY --env production`,
	Args: cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]

		// --force flag skips confirmation prompt
		// Check if user passed --force or -f
		// Like: if (!options.force) { askForConfirmation() }
		force, _ := cmd.Flags().GetBool("force")

		// Ask for confirmation before deleting
		// Deleting a secret is destructive — we want to be sure
		// Skip this if --force flag was passed
		if !force {
			fmt.Printf("⚠️  Are you sure you want to remove %s? (y/N): ", key)

			// Read user input from terminal
			// bufio.NewReader wraps os.Stdin to read line by line
			// Like: readline in Node.js
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n') // read until Enter key
			if err != nil {
				return err
			}

			// Clean up the response
			// strings.TrimSpace removes spaces and newlines
			// strings.ToLower makes it case insensitive (Y and y both work)
			response = strings.TrimSpace(strings.ToLower(response))

			// If user did not type "y" or "yes" → cancel
			if response != "y" && response != "yes" {
				fmt.Println("Cancelled.")
				return nil
			}
		}

		// Now ask for passphrase
		passphrase, err := promptPassphrase()
		if err != nil {
			return err
		}

		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		v, err := vault.Load(dir, passphrase)
		if err != nil {
			return err
		}

		// Remove the secret
		// envFlag from --env flag
		if err := v.Remove(key, envFlag); err != nil {
			return err
		}

		env := envFlag
		if env == "" {
			env = "all environments"
		}

		fmt.Printf("✅ Removed %s (%s)\n", key, env)
		return nil
	},
}

func init() {
	removeCmd.Flags().StringVarP(&envFlag, "env", "e", "", "environment")

	// BoolP creates a boolean flag
	// --force or -f skips confirmation
	// false = default value (confirmation ON by default)
	removeCmd.Flags().BoolP("force", "f", false, "skip confirmation prompt")

	rootCmd.AddCommand(removeCmd)
}
