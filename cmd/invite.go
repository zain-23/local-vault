package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/zain-23/local-vault/internal/client"
	"github.com/zain-23/local-vault/internal/config"
	"github.com/zain-23/local-vault/internal/identity"
	internalsync "github.com/zain-23/local-vault/internal/sync"
)

var inviteCmd = &cobra.Command{
	Use:   "invite",
	Short: "Generate a join token for a teammate",
	Example: `  lv invite --name "Ahmed"
  lv invite --name "Sara" --expires 7d
  lv invite --list
  lv invite --revoke lv_join_xxx`,
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

		if cfg.VaultID == "" {
			return fmt.Errorf(
				"vault not registered on server\n  Run: lv push first",
			)
		}

		sc := client.New(cfg.SignalingServer, id.DeviceID)
		if err := sc.HealthCheck(); err != nil {
			return err
		}

		// Handle --list flag
		listFlag, _ := cmd.Flags().GetBool("list")
		if listFlag {
			return listTokens(sc, cfg.VaultID)
		}

		// Generate new token
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			return fmt.Errorf(
				"provide a name for the token\n  Example: lv invite --name \"Ahmed\"",
			)
		}

		// The shared vault key (DEK) must travel to the joiner inside the
		// token. Load the vault to read it (generating one for vaults that
		// predate this feature).
		v, err := loadVault(dir)
		if err != nil {
			return err
		}
		dek, err := v.EnsureDataKey()
		if err != nil {
			return err
		}

		// Random per-token secret — never sent to the server in usable form.
		secret := make([]byte, 24)
		if _, err := rand.Read(secret); err != nil {
			return fmt.Errorf("failed to generate token secret: %w", err)
		}

		wrappedDEK, err := internalsync.WrapKey(dek, secret)
		if err != nil {
			return err
		}

		token, err := sc.CreateToken(cfg.VaultID, client.CreateTokenRequest{
			DeviceID:   id.DeviceID,
			Name:       name,
			WrappedDEK: wrappedDEK,
			Verifier:   internalsync.DeriveVerifier(secret),
		})
		if err != nil {
			return fmt.Errorf("failed to create token: %w", err)
		}

		// Full shareable token = public id + "." + secret.
		fullToken := token.ID + "." + base64.RawURLEncoding.EncodeToString(secret)

		// Display token
		fmt.Println()
		fmt.Println("╔══════════════════════════════════════════════════╗")
		fmt.Println("║           🔐 LocalVault Join Token               ║")
		fmt.Println("╠══════════════════════════════════════════════════╣")
		fmt.Printf("║  For      : %-36s║\n", name)
		fmt.Println("╠══════════════════════════════════════════════════╣")
		fmt.Println("║  Share this token privately with your teammate   ║")
		fmt.Println("║  ✅ No expiry — works until revoked              ║")
		fmt.Println("║  ✅ Multiple teammates can use different tokens   ║")
		fmt.Printf("║  ✅ Revoke anytime: lv invite --revoke %-10s║\n", token.ID[:10]+"…")
		fmt.Println("╚══════════════════════════════════════════════════╝")
		fmt.Println()
		fmt.Println("They should run:")
		fmt.Printf("  lv join %s\n", fullToken)
		fmt.Println()
		fmt.Println("⚠️  This token contains the vault key — share it privately.")

		return nil
	},
}

func listTokens(sc *client.Client, vaultID string) error {
	tokens, err := sc.ListTokens(vaultID)
	if err != nil {
		return err
	}

	if len(tokens) == 0 {
		fmt.Println("No active tokens.")
		fmt.Println("Create one: lv invite --name \"Ahmed\"")
		return nil
	}

	fmt.Println("🎟️  Active Join Tokens")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	for _, t := range tokens {
		expires := "never"
		if t.ExpiresAt != nil {
			expires = t.ExpiresAt.Format("2006-01-02")
		}
		fmt.Printf("\n  Name    : %s\n", t.Name)
		fmt.Printf("  Token   : %s\n", t.ID)
		fmt.Printf("  Created : %s\n", t.CreatedAt.Format("2006-01-02 15:04"))
		fmt.Printf("  Expires : %s\n", expires)
	}
	fmt.Printf("\n%d active token(s)\n", len(tokens))
	return nil
}

func init() {
	inviteCmd.Flags().StringP("name", "n", "", "name of the teammate (required)")
	inviteCmd.Flags().Bool("list", false, "list all active tokens")
	rootCmd.AddCommand(inviteCmd)
}
