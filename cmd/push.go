package cmd

import (
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

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push secrets to server and all peers",
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
		if err := sc.HealthCheck(); err != nil {
			return err
		}

		// Get all secrets
		vaultSecrets := v.GetSecretEntries()
		syncSecrets := make([]internalsync.SecretEntry, len(vaultSecrets))
		for i, s := range vaultSecrets {
			syncSecrets[i] = internalsync.SecretEntry{
				Key:       s.Key,
				Value:     s.Value,
				Env:       s.Env,
				UpdatedAt: s.UpdatedAt,
			}
		}

		// Get peers from server (most up to date list)
		serverPeers, err := sc.GetVaultPeers(cfg.VaultID)
		if err != nil {
			// Fall back to local peers
			fmt.Println("⚠️  Using local peer list")
		} else {
			// Sync server peers to local vault
			for _, sp := range serverPeers {
				if sp.DeviceID == id.DeviceID {
					continue
				}
				_, found := v.GetPeer(sp.DeviceID)
				if !found {
					v.AddPeer(vault.Peer{
						DeviceID:        sp.DeviceID,
						DeviceName:      sp.DeviceName,
						PublicKey:       sp.PublicKey,
						X25519PublicKey: sp.X25519PublicKey,
					})
				}
			}
		}

		peers := v.GetPeers()
		if len(peers) == 0 {
			fmt.Println("No peers yet.")
			fmt.Println("Invite teammates: lv invite --name \"Ahmed\"")
			return nil
		}

		fmt.Printf("📤 Pushing %d secret(s) to %d peer(s)...\n",
			len(syncSecrets), len(peers))

		// ── Step 1: Send to each peer directly ─────────────────
		for _, peer := range peers {
			if peer.X25519PublicKey == nil {
				continue
			}

			payload, err := internalsync.EncryptForPeer(
				syncSecrets,
				id.X25519PrivateKey,
				peer.X25519PublicKey,
				id.DeviceID,
			)
			if err != nil {
				fmt.Printf("  ⚠️  Failed to encrypt for %s\n", peer.DeviceName)
				continue
			}

			err = sc.SendMessage(client.SendMessageRequest{
				ForDeviceID:   peer.DeviceID,
				FromDeviceID:  id.DeviceID,
				FromPublicKey: id.X25519PublicKey,
				Payload:       payload,
			})
			if err != nil {
				fmt.Printf("  ⚠️  Failed to send to %s\n", peer.DeviceName)
				continue
			}

			fmt.Printf("  ✅ Sent to %s\n", peer.DeviceName)
		}

		// ── Step 2: Upload snapshot to server ──────────────────
		// Snapshot = secrets encrypted for ANY peer to decrypt on join
		// Use first peer's key for snapshot encryption
		// (anyone with valid token gets decryption ability via join)
		fmt.Println()
		fmt.Println("📡 Uploading snapshot to server...")

		if len(peers) > 0 && peers[0].X25519PublicKey != nil {
			snapshot, err := internalsync.EncryptForPeer(
				syncSecrets,
				id.X25519PrivateKey,
				peers[0].X25519PublicKey,
				id.DeviceID,
			)
			if err == nil {
				if err := sc.UploadSnapshot(cfg.VaultID, snapshot); err != nil {
					fmt.Printf("  ⚠️  Snapshot upload failed: %v\n", err)
				} else {
					fmt.Println("  ✅ Snapshot updated on server")
					fmt.Println("     New joiners will get secrets immediately")
				}
			}
		}

		// ── Step 3: Share peer list for full mesh ───────────────
		for _, peer := range peers {
			if peer.X25519PublicKey == nil {
				continue
			}

			var otherPeers []vault.Peer
			for _, other := range peers {
				if other.DeviceID != peer.DeviceID {
					otherPeers = append(otherPeers, other)
				}
			}

			if len(otherPeers) == 0 {
				continue
			}

			peerJSON, _ := json.Marshal(otherPeers)
			peerPayload := append([]byte("peers:"), peerJSON...)

			sc.SendMessage(client.SendMessageRequest{
				ForDeviceID:   peer.DeviceID,
				FromDeviceID:  id.DeviceID,
				FromPublicKey: id.X25519PublicKey,
				Payload:       peerPayload,
			})
		}

		fmt.Println()
		fmt.Println("✅ Push complete")
		fmt.Println("   Peers receive secrets on next: lv sync")
		fmt.Println("   New joiners get secrets immediately on: lv join")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
