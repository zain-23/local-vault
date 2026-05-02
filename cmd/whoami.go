package cmd

// whoami.go handles "lv whoami"
// Shows this device's identity information
// Like: ssh-keygen -l to show your SSH key fingerprint

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zain-23/local-vault/internal/identity"
)

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show this device's identity",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Find .lv directory
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		// Load identity from .lv/identity.json
		lvDir := dir + "/.lv"
		id, err := identity.Load(lvDir)
		if err != nil {
			return err
		}

		// Display identity info
		fmt.Println("🔑 Device Identity")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Printf("Device Name : %s\n", id.DeviceName)
		fmt.Printf("Device ID   : %s\n", id.DeviceID)
		fmt.Printf("Created At  : %s\n", id.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println()
		fmt.Println("Public Key (safe to share):")
		fmt.Println(id.PublicKeyString())

		return nil
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}
