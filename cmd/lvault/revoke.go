package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zain-23/local-vault/internal/client"
	"github.com/zain-23/local-vault/internal/config"
	"github.com/zain-23/local-vault/internal/identity"
	internalsync "github.com/zain-23/local-vault/internal/sync"
)

var revokeCmd = &cobra.Command{
	Use:   "revoke DEVICE_ID",
	Short: "Remove a peer's access to the vault",
	Example: `  lv revoke b427d8d8-8b7f-4605-9b5a-f8be956b168d

  Tip: run "lv peers" to see all device IDs`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		deviceID := args[0]

		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		lvDir := filepath.Join(dir, ".lv")

		// Load vault
		v, err := loadVault(dir)
		if err != nil {
			return err
		}

		// Check peer exists locally
		peer, found := v.GetPeer(deviceID)
		if !found {
			return fmt.Errorf(
				"peer not found: %s\nRun 'lv peers' to see all peers",
				deviceID,
			)
		}

		// Confirm before revoking
		force, _ := cmd.Flags().GetBool("force")
		if !force {
			fmt.Printf("⚠️  Revoke access for %s?\n", peer.DeviceName)
			fmt.Printf("   Device ID : %s\n", deviceID)
			fmt.Println()
			fmt.Println("   This will:")
			fmt.Println("   → Remove them from local trusted peers")
			fmt.Println("   → Remove them from server vault")
			fmt.Println("   → They will no longer receive secret updates")
			fmt.Println("   → They cannot sync anymore")
			fmt.Println()
			fmt.Print("Continue? (y/N): ")

			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				return err
			}

			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				fmt.Println("Cancelled.")
				return nil
			}
		}

		// Load identity and config
		id, err := identity.Load(lvDir)
		if err != nil {
			return err
		}

		cfg, err := config.Load(lvDir)
		if err != nil {
			return err
		}

		// Create signaling client
		sc := client.New(cfg.SignalingServer, id.DeviceID)

		// Step 1 — Remove peer from LOCAL vault
		if err := v.RemovePeer(deviceID); err != nil {
			return err
		}
		fmt.Printf("✅ Removed %s from local peers\n", peer.DeviceName)

		// Step 2 — Remove peer from SERVER vault
		// Critical — prevents revoked peer from syncing
		if cfg.VaultID != "" {
			if err := sc.HealthCheck(); err != nil {
				fmt.Println("⚠️  Could not reach server")
				fmt.Println("   They may still sync until server is updated")
				fmt.Println("   Try again when online: lv revoke", deviceID)
			} else {
				if err := sc.RemovePeer(cfg.VaultID, deviceID); err != nil {
					fmt.Printf("⚠️  Could not remove from server: %v\n", err)
				} else {
					fmt.Printf("✅ Removed %s from server\n", peer.DeviceName)
				}
			}
		}

		// Step 3 — Notify remaining peers
		remainingPeers := v.GetPeers()
		if len(remainingPeers) == 0 {
			fmt.Println("💡 No other peers to notify")
			printRevokeNextSteps(peer.DeviceName)
			return nil
		}

		// Get updated secrets to push to remaining peers
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

		fmt.Printf("📤 Notifying %d remaining peer(s)...\n",
			len(remainingPeers))

		for _, p := range remainingPeers {
			if p.X25519PublicKey == nil {
				continue
			}

			payload, err := internalsync.EncryptForPeer(
				syncSecrets,
				id.X25519PrivateKey,
				p.X25519PublicKey,
				id.DeviceID,
			)
			if err != nil {
				fmt.Printf("  ⚠️  Failed to encrypt for %s\n", p.DeviceName)
				continue
			}

			err = sc.SendMessage(client.SendMessageRequest{
				ForDeviceID:   p.DeviceID,
				FromDeviceID:  id.DeviceID,
				FromPublicKey: id.X25519PublicKey,
				Payload:       payload,
			})
			if err != nil {
				fmt.Printf("  ⚠️  Failed to notify %s\n", p.DeviceName)
				continue
			}

			fmt.Printf("  ✅ Notified %s\n", p.DeviceName)
		}

		printRevokeNextSteps(peer.DeviceName)
		return nil
	},
}

// printRevokeNextSteps shows what to do after revoking
func printRevokeNextSteps(name string) {
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("✅ %s has been revoked\n", name)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("⚠️  Important next steps:")
	fmt.Println("   1. lv rotate --all")
	fmt.Println("      Invalidates their cached secrets")
	fmt.Println()
	fmt.Println("   2. lv push")
	fmt.Println("      Updates remaining peers with rotated secrets")
}

func init() {
	revokeCmd.Flags().BoolP("force", "f", false, "skip confirmation prompt")
	rootCmd.AddCommand(revokeCmd)
}
