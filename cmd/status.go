package cmd

// status.go handles "lv status"
// Shows vault health — secrets count, peers, last sync
// Like "git status" but for your vault

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/zain-23/local-vault/internal/identity"
	"github.com/zain-23/local-vault/internal/vault"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show vault status and sync info",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		lvDir := filepath.Join(dir, ".lv")

		// Load identity — no passphrase needed for this
		id, err := identity.Load(lvDir)
		if err != nil {
			return err
		}

		// Load vault
		passphrase, err := promptPassphrase()
		if err != nil {
			return err
		}

		v, err := vault.Load(dir, passphrase)
		if err != nil {
			return err
		}

		secrets := v.List("")
		peers := v.GetPeers()

		// Count secrets per environment
		envCounts := map[string]int{}
		for _, s := range secrets {
			env := s.Env
			if env == "" {
				env = "all"
			}
			envCounts[env]++
		}

		fmt.Println("🔐 LocalVault Status")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Printf("  Device    : %s\n", id.DeviceName)
		fmt.Printf("  Device ID : %s\n", id.DeviceID)
		fmt.Println()
		fmt.Printf("  Secrets   : %d total\n", len(secrets))

		// Show per environment breakdown
		for env, count := range envCounts {
			fmt.Printf("    ├─ %s: %d\n", env, count)
		}

		fmt.Println()
		fmt.Printf("  Peers     : %d trusted\n", len(peers))
		for _, peer := range peers {
			fmt.Printf("    ├─ %s (%s)\n", peer.DeviceName, peer.DeviceID[:8]+"...")
		}

		fmt.Println()
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println("  Run: lv sync   to pull latest secrets")
		fmt.Println("  Run: lv push   to send secrets to peers")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
