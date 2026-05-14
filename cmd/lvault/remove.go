package main

// remove.go handles "lv remove KEY"
// Deletes a secret from the vault permanently

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove KEY",
	Short: "Remove a secret from the vault",
	Example: `  lv remove DATABASE_URL
  lv remove STRIPE_KEY --env production`,
	Args: cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]

		force, _ := cmd.Flags().GetBool("force")

		// Confirm before deleting — destructive action
		if !force {
			fmt.Printf("⚠️  Are you sure you want to remove %s? (y/N): ", key)

			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				return err
			}

			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				fmt.Println("Cancelled.")
				return nil
			}
		}

		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		// Use session key — no passphrase needed
		v, err := loadVault(dir)
		if err != nil {
			return err
		}

		if err := v.Remove(key, envFlag); err != nil {
			return err
		}

		env := envFlag
		if env == "" {
			env = "all environments"
		}

		fmt.Printf("✅ Removed %s (%s)\n", key, env)
		fmt.Println("💡 Run: lv push   to sync with peers")
		return nil
	},
}

func init() {
	removeCmd.Flags().StringVarP(&envFlag, "env", "e", "", "environment")
	removeCmd.Flags().BoolP("force", "f", false, "skip confirmation prompt")
	rootCmd.AddCommand(removeCmd)
}
