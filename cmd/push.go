package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/zain-23/local-vault/internal/client"
	"github.com/zain-23/local-vault/internal/config"
	"github.com/zain-23/local-vault/internal/identity"
	internalsync "github.com/zain-23/local-vault/internal/sync"
	"github.com/zain-23/local-vault/internal/vault"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push your secrets to all trusted peers",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		lvDir := filepath.Join(dir, ".lv")

		// Load identity
		id, err := identity.Load(lvDir)
		if err != nil {
			return err
		}

		// Load config
		cfg, err := config.Load(lvDir)
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

		// Get trusted peers
		peers := v.GetPeers()
		if len(peers) == 0 {
			fmt.Println("No trusted peers yet.")
			fmt.Println("Share an invite: lv invite")
			return nil
		}

		// Create signaling client
		sc := client.New(cfg.SignalingServer, id.DeviceID)

		if err := sc.HealthCheck(); err != nil {
			return err
		}

		// Get secrets in transfer format
		vaultSecrets := v.GetSecretEntries()

		// Convert vault.SecretEntry to sync.SecretEntry
		syncSecrets := make([]internalsync.SecretEntry, len(vaultSecrets))
		for i, s := range vaultSecrets {
			syncSecrets[i] = internalsync.SecretEntry{
				Key:       s.Key,
				Value:     s.Value,
				Env:       s.Env,
				UpdatedAt: s.UpdatedAt,
			}
		}

		fmt.Printf("📤 Pushing %d secret(s) to %d peer(s)...\n",
			len(syncSecrets), len(peers))

		// Encrypt and send to each peer
		for _, peer := range peers {
			// Need peer's X25519 public key to encrypt for them
			if peer.X25519PublicKey == nil {
				fmt.Printf("  ⚠️  No encryption key for %s — skipping\n",
					peer.DeviceName)
				continue
			}

			// Encrypt using our X25519 private + peer's X25519 public
			// Only this peer can decrypt — their private key required
			payload, err := internalsync.EncryptForPeer(
				syncSecrets,
				id.X25519PrivateKey,  // our private key
				peer.X25519PublicKey, // their public key
				id.DeviceID,
			)
			if err != nil {
				fmt.Printf("  ⚠️  Failed to encrypt for %s: %v\n",
					peer.DeviceName, err)
				continue
			}

			// Send to peer's mailbox on signaling server
			// If peer online → they get it immediately
			// If peer offline → stored for 48 hours
			err = sc.SendMessage(client.SendMessageRequest{
				ForDeviceID:   peer.DeviceID,
				FromDeviceID:  id.DeviceID,
				FromPublicKey: id.X25519PublicKey,
				Payload:       payload,
			})
			if err != nil {
				fmt.Printf("  ⚠️  Failed to send to %s: %v\n",
					peer.DeviceName, err)
				continue
			}

			fmt.Printf("  ✅ Sent to %s\n", peer.DeviceName)
		}

		fmt.Println()
		fmt.Println("✅ Push complete")
		fmt.Println("   Peers will receive secrets on next: lv sync")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}

// Ensure vault import used
var _ = vault.Peer{}
