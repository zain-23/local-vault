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

var joinCmd = &cobra.Command{
	Use:     "join CODE",
	Short:   "Join a teammate's vault using their invite code",
	Example: "  lv join LV-A3F9-X2K1",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		code := args[0]

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

		// Create signaling client
		sc := client.New(cfg.SignalingServer, id.DeviceID)

		fmt.Println("🔍 Looking up invite code...")
		if err := sc.HealthCheck(); err != nil {
			return err
		}

		// Look up invite — returns Dev A's device ID and X25519 public key
		peer, err := sc.LookupInvite(code)
		if err != nil {
			return err
		}

		fmt.Printf("✅ Found peer: %s\n", peer.DeviceID)
		fmt.Println("🤝 Exchanging keys...")

		// Load vault using session — no passphrase needed
		v, err := loadVault(dir)
		if err != nil {
			return err
		}

		// Save Dev A as trusted peer
		err = v.AddPeer(vault.Peer{
			DeviceID:        peer.DeviceID,
			DeviceName:      peer.DeviceID,
			PublicKey:       peer.PublicKey,
			X25519PublicKey: peer.PublicKey,
		})
		if err != nil {
			return fmt.Errorf("failed to save peer: %w", err)
		}

		fmt.Println("📤 Sending your public key to peer...")

		// Send hello to Dev A's mailbox with our X25519 public key
		err = sc.SendMessage(client.SendMessageRequest{
			ForDeviceID:   peer.DeviceID,
			FromDeviceID:  id.DeviceID,
			FromPublicKey: id.X25519PublicKey,
			Payload:       []byte("hello"),
		})
		if err != nil {
			return fmt.Errorf("failed to notify peer: %w", err)
		}

		fmt.Println("📬 Checking for vault data...")

		// Check if Dev A already sent us secrets
		msgs, err := sc.GetMessages()
		if err != nil {
			return err
		}

		if msgs.Count == 0 {
			fmt.Println()
			fmt.Println("⏳ Your key was sent to peer.")
			fmt.Println("   Ask them to run: lv push")
			fmt.Println("   Then run: lv sync")
			return nil
		}

		// Process messages already waiting
		totalMerged := 0
		for _, msg := range msgs.Messages {
			// Skip hello messages
			if string(msg.Payload) == "hello" {
				continue
			}

			// Decrypt using our X25519 private + peer's X25519 public
			rawSecrets, err := internalsync.DecryptFromPeer(
				msg.Payload,
				id.X25519PrivateKey,
				peer.PublicKey,
			)
			if err != nil {
				fmt.Printf("⚠️  Could not decrypt: %v\n", err)
				continue
			}

			// Convert to vault.SecretEntry
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
		}

		fmt.Println()
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Printf("✅ Connected to peer: %s\n", peer.DeviceID)
		if totalMerged > 0 {
			fmt.Printf("📦 Synced %d secret(s)\n", totalMerged)
		} else {
			fmt.Println("📦 No secrets yet — ask peer to run: lv push")
		}
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println()
		fmt.Println("Run: lv sync   to pull secrets anytime")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(joinCmd)
}
