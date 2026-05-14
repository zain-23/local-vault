package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add KEY=VALUE",
	Short: "Add or update a secret",
	Example: `  lv add DATABASE_URL=postgres://localhost/mydb
  lv add API_KEY=sk-xxx --env production
  lv add SECRET_KEY=abc123 --env staging`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Parse KEY=VALUE
		parts := strings.SplitN(args[0], "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid format. use: lv add KEY=VALUE")
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "" {
			return fmt.Errorf("key cannot be empty")
		}

		dir, _ := os.Getwd()

		// Use session key — no passphrase needed
		v, err := loadVault(dir)
		if err != nil {
			return err
		}

		if err := v.Add(key, value, envFlag); err != nil {
			return err
		}

		env := envFlag
		if env == "" {
			env = "all environments"
		}

		fmt.Printf("✅ Added %s (%s)\n", key, env)
		fmt.Println("💡 Run: lv push   to sync with peers")
		return nil
	},
}

func init() {
	addCmd.Flags().StringVarP(&envFlag, "env", "e", "", "environment (development/staging/production)")
	rootCmd.AddCommand(addCmd)
}
