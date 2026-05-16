package cmd

// log.go handles "lv log"
// Shows audit trail of all changes made to the vault
// Like "git log" but for secrets

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show audit trail of vault changes",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		// Use session key — no passphrase needed
		v, err := loadVault(dir)
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
			fmt.Printf("   Device : %s\n", shortID(entry.DeviceID))
		}

		fmt.Printf("\n%d entries total\n", len(entries))
		return nil
	},
}

// shortID truncates a device identifier for display.
// Peer device IDs are UUIDs (36 chars); locally-originated audit
// entries use the short sentinel "local". Only truncate when there
// is actually something to hide, otherwise return the value as-is.
func shortID(id string) string {
	if len(id) <= 8 {
		return id
	}
	return id[:8] + "..."
}

func init() {
	rootCmd.AddCommand(logCmd)
}
