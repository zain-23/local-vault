package cmd

// log.go handles "lv log"
// Shows audit trail of all changes made to the vault
// Like "git log" but for secrets

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zain-23/local-vault/internal/vault"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show audit trail of vault changes",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		passphrase, err := promptPassphrase()
		if err != nil {
			return err
		}

		v, err := vault.Load(dir, passphrase)
		if err != nil {
			return err
		}

		entries := v.GetAuditLog()

		if len(entries) == 0 {
			fmt.Println("No audit log entries yet.")
			fmt.Println("Entries are recorded when you add, update, or remove secrets.")
			return nil
		}

		fmt.Println("📋 Vault Audit Log")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

		for _, entry := range entries {
			// Format action with emoji
			var emoji string
			switch entry.Action {
			case "add":
				emoji = "➕"
			case "update":
				emoji = "✏️ "
			case "remove":
				emoji = "🗑️ "
			case "rotate":
				emoji = "🔄"
			default:
				emoji = "•"
			}

			env := entry.Env
			if env == "" {
				env = "all"
			}

			fmt.Printf("\n%s  %s\n", emoji, entry.Action)
			fmt.Printf("   Key    : %s [%s]\n", entry.Key, env)
			fmt.Printf("   When   : %s\n",
				entry.Timestamp.Format("2006-01-02 15:04:05"))
			fmt.Printf("   Device : %s\n", entry.DeviceID[:8]+"...")
		}

		fmt.Printf("\n%d entries total\n", len(entries))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
}
