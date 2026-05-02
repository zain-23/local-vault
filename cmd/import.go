package cmd

// import.go handles "lv import .env.local"
// Reads an existing .env file and adds all secrets to vault
// This is how teams migrate FROM .env files TO LocalVault

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zain-23/local-vault/internal/vault"
)

var importCmd = &cobra.Command{
	Use:   "import FILE",
	Short: "Import secrets from a .env file into vault",
	Example: `  lv import .env.local
  lv import .env.production --env production
  lv import .env.development --env development`,

	// ExactArgs(1) = user must provide the file path
	Args: cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		// args[0] = file path user provided
		// Example: ".env.local" or "/home/user/project/.env"
		filePath := args[0]

		// Check if file exists BEFORE asking for passphrase
		// Better UX — fail fast if file not found
		// Like: fs.existsSync(filePath) in Node.js
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", filePath)
		}

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

		// ImportEnvFile is defined in internal/vault/vault.go
		// Reads the file, parses KEY=VALUE lines, adds to vault
		// Returns count of how many secrets were imported
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
		fmt.Printf("   Add it to .gitignore if not already there\n")

		// Show what was imported
		// List all secrets so user can verify
		fmt.Println()
		fmt.Println("Imported secrets:")
		secrets := v.List(envFlag)
		for _, s := range secrets {
			// Print key name, mask value with asterisks
			fmt.Printf("  ✓ %s\n", s.Key)
		}

		return nil
	},
}

func init() {
	importCmd.Flags().StringVarP(&envFlag, "env", "e", "", "environment to import into")
	rootCmd.AddCommand(importCmd)
}
