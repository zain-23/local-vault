package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/zain-23/local-vault/internal/api"
	"github.com/zain-23/local-vault/internal/identity"
	internalsync "github.com/zain-23/local-vault/internal/sync"
	"github.com/zain-23/local-vault/internal/ui"
)

var inviteCmd = &cobra.Command{
	Use:   "invite",
	Short: "Generate a join token for a teammate",
	Example: `  lv invite --name "Ahmed"
  lv invite --list`,
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
		cfg, err := requireLinkedConfig(lvDir)
		if err != nil {
			return err
		}
		client, err := requireAPI()
		if err != nil {
			return err
		}

		listFlag, _ := cmd.Flags().GetBool("list")
		if listFlag {
			return listTokens(client, cfg.WorkspaceID, cfg.VaultID)
		}

		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			return fmt.Errorf("provide a name for the token\n  Example: lv invite --name \"Ahmed\"")
		}

		v, err := loadVault(dir)
		if err != nil {
			return err
		}
		dek, err := v.EnsureDataKey()
		if err != nil {
			return err
		}

		secret := make([]byte, 24)
		if _, err := rand.Read(secret); err != nil {
			return fmt.Errorf("failed to generate token secret: %w", err)
		}
		wrappedDEK, err := internalsync.WrapKey(dek, secret)
		if err != nil {
			return err
		}

		token, err := client.CreateToken(cfg.WorkspaceID, cfg.VaultID, api.CreateTokenRequest{
			DeviceID:   id.DeviceID,
			Name:       name,
			WrappedDEK: wrappedDEK,
			Verifier:   internalsync.DeriveVerifier(secret),
		})
		if err != nil {
			return mapNotLoggedIn(fmt.Errorf("failed to create token: %w", err))
		}

		fullToken := token.ID + "." + base64.RawURLEncoding.EncodeToString(secret)

		ui.Header("Join Token")
		ui.KeyValue("For", name)
		ui.KeyValue("Token ID", token.ID)
		ui.Info("share this privately with your teammate:")
		ui.Code(fullToken)
		ui.Hint("they should run: lv join <token>")
		ui.Warn("this token contains the vault key — share it privately")
		return nil
	},
}

func listTokens(client *api.Client, workspaceID, vaultID string) error {
	tokens, err := client.ListTokens(workspaceID, vaultID)
	if err != nil {
		return mapNotLoggedIn(err)
	}
	if len(tokens) == 0 {
		ui.Info("no active tokens")
		ui.Hint("create one: lv invite --name \"Ahmed\"")
		return nil
	}
	rows := make([][]string, 0, len(tokens))
	for _, t := range tokens {
		expires := "never"
		if t.ExpiresAt != nil {
			expires = t.ExpiresAt.Format("2006-01-02")
		}
		rows = append(rows, []string{t.Name, t.ID, t.CreatedAt.Format("2006-01-02 15:04"), expires})
	}
	ui.Header("Active Join Tokens")
	ui.Table([]string{"NAME", "TOKEN", "CREATED", "EXPIRES"}, rows)
	return nil
}

func init() {
	inviteCmd.Flags().StringP("name", "n", "", "name of the teammate (required)")
	inviteCmd.Flags().Bool("list", false, "list all active tokens")
	rootCmd.AddCommand(inviteCmd)
}
