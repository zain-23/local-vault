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

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Pull latest secrets from peers",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		lvDir := filepath.Join(dir, ".lv")

		// Load identity — no passphrase needed
		id, err := identity.Load(lvDir)
		if err != nil {
			return err
		}

		// Load config
		cfg, err := config.Load(lvDir)
		if err != nil {
			return err
		}

		// Load vault using session — no passphrase needed
		v, err := loadVault(dir)
		if err != nil {
			return err
		}

		// Create signaling client
		sc := client.New(cfg.SignalingServer, id.DeviceID)

		fmt.Println("🔄 Syncing with peers...")

		if err := sc.HealthCheck(); err != nil {
			return err
		}

		// Check mailbox for pending messages
		msgs, err := sc.GetMessages()
		if err != nil {
			return err
		}

		if msgs.Count == 0 {
			fmt.Println("✅ Already up to date")
			return nil
		}

		fmt.Printf("📬 Found %d update(s)\n", msgs.Count)

		totalMerged := 0
		for _, msg := range msgs.Messages {
			// Find sender in trusted peers
			peer, found := v.GetPeer(msg.FromDeviceID)
			if !found {
				// Unknown peer — save if they sent public key
				if msg.FromPublicKey != nil {
					err = v.AddPeer(vault.Peer{
						DeviceID:        msg.FromDeviceID,
						DeviceName:      msg.FromDeviceID,
						PublicKey:       msg.FromPublicKey,
						X25519PublicKey: msg.FromPublicKey,
					})
					if err != nil {
						fmt.Printf("⚠️  Could not save peer: %v\n", err)
						continue
					}
					peer, _ = v.GetPeer(msg.FromDeviceID)
					fmt.Printf("✅ New peer discovered: %s\n", msg.FromDeviceID)
				} else {
					fmt.Printf("⚠️  Message from unknown peer %s — skipping\n",
						msg.FromDeviceID)
					continue
				}
			}

			// Skip hello messages — peer discovery only
			if string(msg.Payload) == "hello" {
				fmt.Printf("👋 Hello from %s — peer saved\n", msg.FromDeviceID)
				continue
			}

			// Check peer has X25519 key
			if peer.X25519PublicKey == nil {
				fmt.Printf("⚠️  No encryption key for peer %s — skipping\n",
					msg.FromDeviceID)
				continue
			}

			// Decrypt using our X25519 private + sender's X25519 public
			rawSecrets, err := internalsync.DecryptFromPeer(
				msg.Payload,
				id.X25519PrivateKey,
				peer.X25519PublicKey,
			)
			if err != nil {
				fmt.Printf("⚠️  Could not decrypt message from %s: %v\n",
					msg.FromDeviceID, err)
				continue
			}

			if len(rawSecrets) == 0 {
				continue
			}

			// Convert sync.SecretEntry to vault.SecretEntry
			secrets := make([]vault.SecretEntry, len(rawSecrets))
			for i, s := range rawSecrets {
				secrets[i] = vault.SecretEntry{
					Key:       s.Key,
					Value:     s.Value,
					Env:       s.Env,
					UpdatedAt: s.UpdatedAt,
				}
			}

			// Merge into vault — newer timestamp wins
			count, err := v.MergeSecrets(secrets)
			if err != nil {
				return err
			}

			totalMerged += count
			fmt.Printf("  ✅ Received %d secret(s) from %s\n",
				count, msg.FromDeviceID)
		}

		fmt.Println()
		if totalMerged > 0 {
			fmt.Printf("✅ Synced %d secret(s) total\n", totalMerged)
		} else {
			fmt.Println("✅ No new changes")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
