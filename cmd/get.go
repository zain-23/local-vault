package cmd

// get.go handles the "lv get KEY" command
// Shows the actual value of a single secret
// Unlike list which shows all keys but hides values

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zain-23/local-vault/internal/vault"
)

var getCmd = &cobra.Command{
	Use:   "get KEY",
	Short: "Get the value of a secret",

	// Example shown when user runs: lv get --help
	Example: `  lv get DATABASE_URL
  lv get STRIPE_KEY --env production
  lv get API_KEY --env development`,

	// ExactArgs(1) means user MUST provide exactly 1 argument
	// If they type just "lv get" with no key → shows error automatically
	// Like checking: if (args.length !== 1) throw new Error()
	Args: cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		// args[0] is the key name user typed
		// Example: if user ran "lv get DATABASE_URL"
		// then args[0] = "DATABASE_URL"
		key := args[0]

		// Ask for passphrase to decrypt vault
		passphrase, err := promptPassphrase()
		if err != nil {
			return err
		}

		// Get current directory
		// Vault is always in .lv/ inside current project folder
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		// Load and decrypt vault from disk
		v, err := vault.Load(dir, passphrase)
		if err != nil {
			return err
		}

		// Get the secret value
		// envFlag comes from --env flag (defined below in init())
		// If no --env flag given, envFlag = "" which means search all envs
		value, err := v.Get(key, envFlag)
		if err != nil {
			return err
		}

		// Print just the value — clean output
		// This allows piping: lv get DATABASE_URL | pbcopy
		// (copies value to clipboard on Mac)
		fmt.Fprintln(os.Stdout, value)
		return nil
	},
}

func init() {
	// Register --env flag on this command
	// StringVarP means:
	// &envFlag    → store value in envFlag variable
	// "env"       → long flag name (--env)
	// "e"         → short flag name (-e)
	// ""          → default value (empty = all environments)
	// "..."       → description shown in --help
	getCmd.Flags().StringVarP(&envFlag, "env", "e", "", "environment (development/staging/production)")

	// Add this command to root so "lv get" works
	rootCmd.AddCommand(getCmd)
}
