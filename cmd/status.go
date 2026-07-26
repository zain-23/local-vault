package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/zain-23/local-vault/internal/identity"
	"github.com/zain-23/local-vault/internal/session"
	"github.com/zain-23/local-vault/internal/ui"
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

		id, err := identity.Load(lvDir)
		if err != nil {
			return err
		}
		v, err := loadVault(dir)
		if err != nil {
			return err
		}

		secrets := v.List("")
		peers := v.GetPeers()

		envCounts := map[string]int{}
		for _, s := range secrets {
			env := s.Env
			if env == "" {
				env = "all"
			}
			envCounts[env]++
		}

		lockStatus := "Locked"
		if remaining, err := session.TimeRemaining(lvDir); err == nil {
			hours := int(remaining.Hours())
			minutes := int(remaining.Minutes()) % 60
			lockStatus = fmt.Sprintf("Unlocked (%dh %dm remaining)", hours, minutes)
		}

		ui.Header("LocalVault Status")
		ui.KeyValue("Device", id.DeviceName)
		ui.KeyValue("Device ID", id.DeviceID)
		ui.KeyValue("Session", lockStatus)
		ui.KeyValue("Secrets", fmt.Sprintf("%d total", len(secrets)))
		for env, count := range envCounts {
			ui.KeyValue("  "+env, fmt.Sprintf("%d", count))
		}
		ui.KeyValue("Peers", fmt.Sprintf("%d trusted", len(peers)))
		for _, peer := range peers {
			ui.Info("  %s (%s)", peer.DeviceName, shortID(peer.DeviceID))
		}
		ui.Hint("lv sync   to pull latest secrets")
		ui.Hint("lv push   to send secrets to peers")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
