package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/zain-23/local-vault/internal/ui"
)

var peersCmd = &cobra.Command{
	Use:   "peers",
	Short: "List all trusted peers",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		v, err := loadVault(dir)
		if err != nil {
			return err
		}

		peers := v.GetPeers()
		if len(peers) == 0 {
			ui.Info("no trusted peers yet")
			ui.Hint("share an invite: lv invite")
			return nil
		}

		ui.Header("Trusted Peers")
		for i, peer := range peers {
			ui.Info("Peer %d", i+1)
			ui.KeyValue("Name", peer.DeviceName)
			ui.KeyValue("Device ID", peer.DeviceID)
			ui.KeyValue("Added", peer.AddedAt.Format("2006-01-02 15:04:05"))
		}
		ui.Info("%d peer(s) total", len(peers))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(peersCmd)
}
