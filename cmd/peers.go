package cmd

// peers.go handles "lv peers"
// Shows all trusted peers this vault is connected to
// Useful for debugging invite/sync issues

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var peersCmd = &cobra.Command{
	Use:   "peers",
	Short: "List all trusted peers",
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

		peers := v.GetPeers()

		if len(peers) == 0 {
			fmt.Println("No trusted peers yet.")
			fmt.Println("Share an invite: lv invite")
			return nil
		}

		fmt.Println("👥 Trusted Peers")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

		for i, peer := range peers {
			fmt.Printf("\n  Peer %d\n", i+1)
			fmt.Printf("  Name      : %s\n", peer.DeviceName)
			fmt.Printf("  Device ID : %s\n", peer.DeviceID)
			fmt.Printf("  Added     : %s\n",
				peer.AddedAt.Format("2006-01-02 15:04:05"))
		}

		fmt.Printf("\n%d peer(s) total\n", len(peers))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(peersCmd)
}
