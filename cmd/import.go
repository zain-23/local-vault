package cmd

// import.go handles "lv import .env.local"
// Reads an existing .env file and adds all secrets to vault
// This is how teams migrate FROM .env files TO LocalVault

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import FILE",
	Short: "Import secrets from a .env file into vault",
	Example: `  lv import .env.local
  lv import .env.production --env production
  lv import .env.development --env development`,
	Args: cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]

		// Check file exists BEFORE loading vault
		// Fail fast — better UX
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", filePath)
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

		// Read and parse .env file
		// Adds all KEY=VALUE lines to vault
		count, err := v.ImportEnvFile(filePath, envFlag)
		if err != nil {
			return err
		}

		env := envFlag
		if env == "" {
			env = "all environments"
		}

		fmt.Printf("✅ Imported %d secrets (%s)\n", count, env)
		fmt.Println()
		fmt.Printf("💡 You can now safely delete %s\n", filePath)
		fmt.Println("   Add it to .gitignore if not already there")
		fmt.Println()

		// Show imported keys so user can verify
		fmt.Println("Imported secrets:")
		secrets := v.List(envFlag)
		for _, s := range secrets {
			fmt.Printf("  ✓ %s\n", s.Key)
		}

		fmt.Println()
		fmt.Println("Run: lv push   to sync with peers")

		return nil
	},
}

func init() {
	importCmd.Flags().StringVarP(&envFlag, "env", "e", "", "environment to import into")
	rootCmd.AddCommand(importCmd)
}
