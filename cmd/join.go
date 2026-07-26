package cmd

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/zain-23/local-vault/internal/api"
	"github.com/zain-23/local-vault/internal/appstate"
	"github.com/zain-23/local-vault/internal/authstore"
	"github.com/zain-23/local-vault/internal/config"
	"github.com/zain-23/local-vault/internal/identity"
	internalsync "github.com/zain-23/local-vault/internal/sync"
	"github.com/zain-23/local-vault/internal/ui"
	"github.com/zain-23/local-vault/internal/vault"
)

var joinCmd = &cobra.Command{
	Use:     "join TOKEN",
	Short:   "Join a vault using a join token",
	Example: "  lv join lv_join_a3f9b2c1xxx.<secret>",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		parts := strings.SplitN(args[0], ".", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return fmt.Errorf("invalid token format — expected lv_join_<id>.<secret>")
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

		st, err := appstate.Load()
		if err != nil {
			return err
		}
		client := api.New(st.ServerURL)

		ui.Step("verifying token...")
		resp, err := client.Join(api.JoinRequest{
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

		cfg.VaultID = resp.VaultID
		cfg.WorkspaceID = resp.WorkspaceID
		cfg.DeviceID = id.DeviceID
		if err := config.Save(lvDir, cfg); err != nil {
			return err
		}

		ui.Success("joined vault %s", resp.VaultID)
		ui.Step("discovering %d peer(s)...", len(resp.Peers))
		for _, peer := range resp.Peers {
			if peer.DeviceID == id.DeviceID {
				continue
			}
			if err := v.AddPeer(apiPeerToVault(peer)); err == nil {
				ui.Success("peer saved: %s", peer.DeviceName)
			}
		}

		if resp.WrappedDEK == nil {
			ui.Warn("this invite predates encrypted-key support")
			ui.Hint("ask the owner for a new invite: lv invite --name ...")
			return nil
		}

		dek, err := internalsync.UnwrapKey(resp.WrappedDEK, secret)
		if err != nil {
			return fmt.Errorf("could not unwrap vault key — ask for a new invite")
		}
		if err := v.SetDataKey(dek); err != nil {
			return err
		}

		if resp.Snapshot == nil {
			ui.Warn("no secrets pushed yet — they will arrive on: lv sync")
		} else {
			ui.Step("downloading vault snapshot...")
			rawSecrets, err := internalsync.DecryptSnapshot(resp.Snapshot, dek)
			if err != nil {
				return fmt.Errorf("snapshot decryption failed — vault key mismatch")
			}
			secrets := make([]vault.SecretEntry, len(rawSecrets))
			for i, s := range rawSecrets {
				secrets[i] = vault.SecretEntry{
					Key: s.Key, Value: s.Value, Env: s.Env, UpdatedAt: s.UpdatedAt,
				}
			}
			count, err := v.MergeSecrets(secrets)
			if err != nil {
				return err
			}
			ui.Success("received %d secret(s)", count)
		}

		ui.Header("Joined")
		ui.KeyValue("Vault", resp.VaultID)
		ui.KeyValue("Workspace", resp.WorkspaceID)
		ui.KeyValue("Peers", fmt.Sprintf("%d", len(resp.Peers)))
		ui.Hint("lv inject -- npm run dev")

		if _, err := authstore.Load(); err != nil {
			ui.Hint("run: lv login   before push/sync")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(joinCmd)
}
