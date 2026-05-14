package main

import (
	"bytes"
	"encoding/json"
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
	Short: "Pull latest secrets from server and peers",
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

		cfg, err := config.Load(lvDir)
		if err != nil {
			return err
		}

		v, err := loadVault(dir)
		if err != nil {
			return err
		}

		sc := client.New(cfg.SignalingServer, id.DeviceID)

		fmt.Println("🔄 Syncing...")

		if err := sc.HealthCheck(); err != nil {
			return err
		}

		totalMerged := 0

		// ── Step 1: Sync peer list from server ─────────────────
		// Always get latest peers from server
		// This auto-discovers anyone who joined
		if cfg.VaultID != "" {
			serverPeers, err := sc.GetVaultPeers(cfg.VaultID)
			if err == nil {
				newPeers := 0
				for _, sp := range serverPeers {
					if sp.DeviceID == id.DeviceID {
						continue // skip ourselves
					}
					_, found := v.GetPeer(sp.DeviceID)
					if !found {
						v.AddPeer(vault.Peer{
							DeviceID:        sp.DeviceID,
							DeviceName:      sp.DeviceName,
							PublicKey:       sp.PublicKey,
							X25519PublicKey: sp.X25519PublicKey,
						})
						newPeers++
						fmt.Printf("  🔗 New peer: %s\n", sp.DeviceName)
					}
				}
				if newPeers > 0 {
					fmt.Printf("  ✅ Added %d new peer(s)\n", newPeers)
				}
			}
		}

		// ── Step 2: Check mailbox for messages ─────────────────
		msgs, err := sc.GetMessages()
		if err != nil {
			return err
		}

		if msgs.Count == 0 {
			fmt.Println("✅ Already up to date")
			return nil
		}

		fmt.Printf("📬 Found %d message(s)\n", msgs.Count)

		for _, msg := range msgs.Messages {
			// Find sender in trusted peers
			peer, found := v.GetPeer(msg.FromDeviceID)
			if !found {
				if msg.FromPublicKey != nil {
					v.AddPeer(vault.Peer{
						DeviceID:        msg.FromDeviceID,
						DeviceName:      msg.FromDeviceID,
						PublicKey:       msg.FromPublicKey,
						X25519PublicKey: msg.FromPublicKey,
					})
					peer, _ = v.GetPeer(msg.FromDeviceID)
					fmt.Printf("  ✅ New peer: %s\n", msg.FromDeviceID[:8]+"...")
				} else {
					continue
				}
			}

			// Handle hello messages
			if string(msg.Payload) == "hello" {
				fmt.Printf("  👋 Hello from %s\n", peer.DeviceName)
				continue
			}

			// Handle peer list messages — full mesh
			if bytes.HasPrefix(msg.Payload, []byte("peers:")) {
				peerJSON := msg.Payload[6:]
				var newPeers []vault.Peer
				if err := json.Unmarshal(peerJSON, &newPeers); err != nil {
					continue
				}
				discovered := 0
				for _, np := range newPeers {
					_, exists := v.GetPeer(np.DeviceID)
					if !exists {
						v.AddPeer(np)
						discovered++
						fmt.Printf("  🔗 Discovered: %s\n", np.DeviceName)
					}
				}
				if discovered > 0 {
					fmt.Printf("  ✅ Added %d peer(s) to mesh\n", discovered)
				}
				continue
			}

			// Handle secrets payload
			if peer.X25519PublicKey == nil {
				continue
			}

			rawSecrets, err := internalsync.DecryptFromPeer(
				msg.Payload,
				id.X25519PrivateKey,
				peer.X25519PublicKey,
			)
			if err != nil {
				fmt.Printf("  ⚠️  Could not decrypt from %s\n", peer.DeviceName)
				continue
			}

			if len(rawSecrets) == 0 {
				continue
			}

			secrets := make([]vault.SecretEntry, len(rawSecrets))
			for i, s := range rawSecrets {
				secrets[i] = vault.SecretEntry{
					Key:       s.Key,
					Value:     s.Value,
					Env:       s.Env,
					UpdatedAt: s.UpdatedAt,
				}
			}

			count, err := v.MergeSecrets(secrets)
			if err != nil {
				return err
			}

			totalMerged += count
			fmt.Printf("  ✅ %d secret(s) from %s\n",
				count, peer.DeviceName)
		}

		fmt.Println()
		if totalMerged > 0 {
			fmt.Printf("✅ Synced %d secret(s)\n", totalMerged)
		} else {
			fmt.Println("✅ No new changes")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
