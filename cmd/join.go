package cmd

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zain-23/local-vault/internal/client"
	"github.com/zain-23/local-vault/internal/config"
	"github.com/zain-23/local-vault/internal/identity"
	internalsync "github.com/zain-23/local-vault/internal/sync"
	"github.com/zain-23/local-vault/internal/vault"
)

var joinCmd = &cobra.Command{
	Use:     "join TOKEN",
	Short:   "Join a vault using a join token",
	Example: "  lv join lv_join_a3f9b2c1xxx",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Token = <public id>.<secret>. The secret never goes to the
		// server in usable form; it unwraps the vault key locally.
		parts := strings.SplitN(args[0], ".", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return fmt.Errorf(
				"invalid token format — expected lv_join_<id>.<secret>",
			)
		}
		tokenID := parts[0]
		secret, err := base64.RawURLEncoding.DecodeString(parts[1])
		if err != nil {
			return fmt.Errorf("invalid token format — secret is not valid")
		}

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
		v, err := loadVault(dir)
		if err != nil {
			return err
		}

		// Connect to server
		sc := client.New(cfg.SignalingServer, id.DeviceID)

		fmt.Println("🔍 Verifying token...")
		if err := sc.HealthCheck(); err != nil {
			return err
		}

		// Join vault using token
		// Server returns: vault ID + snapshot + all current peers
		resp, err := sc.JoinVault(client.JoinRequest{
			Token:           tokenID,
			Verifier:        internalsync.DeriveVerifier(secret),
			DeviceID:        id.DeviceID,
			DeviceName:      id.DeviceName,
			PublicKey:       id.PublicKey,
			X25519PublicKey: id.X25519PublicKey,
		})
		if err != nil {
			return fmt.Errorf("failed to join: %w", err)
		}

		fmt.Printf("✅ Joined vault: %s\n", resp.VaultID)

		// Save vault ID to config
		cfg.VaultID = resp.VaultID
		cfg.DeviceID = id.DeviceID
		config.Save(lvDir, cfg)

		// Save ALL peers from server immediately
		// This establishes full mesh from day one
		// No need to wait for lv push from team lead
		fmt.Printf("🔗 Discovering %d peer(s)...\n", len(resp.Peers))
		for _, peer := range resp.Peers {
			// Skip ourselves
			if peer.DeviceID == id.DeviceID {
				continue
			}

			err = v.AddPeer(vault.Peer{
				DeviceID:        peer.DeviceID,
				DeviceName:      peer.DeviceName,
				PublicKey:       peer.PublicKey,
				X25519PublicKey: peer.X25519PublicKey,
			})
			if err == nil {
				fmt.Printf("  ✅ Peer saved: %s\n", peer.DeviceName)
			}
		}

		// Unwrap the shared vault key from the token secret.
		if resp.WrappedDEK == nil {
			fmt.Println()
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			fmt.Printf("✅ Joined vault: %s\n", resp.VaultID)
			fmt.Printf("👥 %d peer(s) in vault\n", len(resp.Peers))
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			fmt.Println()
			fmt.Println("⚠️  This invite predates encrypted-key support —")
			fmt.Println("   you joined as a peer but cannot decrypt secrets.")
			fmt.Println("   Ask the owner for a new invite: lv invite --name ...")
			return nil
		}

		dek, err := internalsync.UnwrapKey(resp.WrappedDEK, secret)
		if err != nil {
			return fmt.Errorf(
				"could not unwrap vault key — token may be from an older version; ask for a new invite",
			)
		}
		if err := v.SetDataKey(dek); err != nil {
			return err
		}

		// Decrypt and merge snapshot if available.
		if resp.Snapshot == nil {
			fmt.Println()
			fmt.Println("⚠️  No secrets pushed yet.")
			fmt.Println("   They will arrive on the next: lv sync")
			fmt.Println()
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			fmt.Printf("✅ Joined vault: %s\n", resp.VaultID)
			fmt.Printf("👥 %d peer(s) in vault\n", len(resp.Peers))
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			return nil
		}

		fmt.Println("📦 Downloading vault snapshot...")

		rawSecrets, err := internalsync.DecryptSnapshot(resp.Snapshot, dek)
		if err != nil {
			return fmt.Errorf("snapshot decryption failed — vault key mismatch")
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
		fmt.Printf("✅ Received %d secret(s)\n", count)

		fmt.Println()
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Printf("✅ Joined vault: %s\n", resp.VaultID)
		fmt.Printf("👥 %d peer(s) in vault\n", len(resp.Peers))
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println()
		fmt.Println("Run: lv inject -- npm run dev")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(joinCmd)
}
