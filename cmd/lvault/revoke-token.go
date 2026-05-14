package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/zain-23/local-vault/internal/client"
	"github.com/zain-23/local-vault/internal/config"
	"github.com/zain-23/local-vault/internal/identity"
)

var revokeTokenCmd = &cobra.Command{
	Use:   "revoke-token TOKEN",
	Short: "Revoke a join token",
	Example: `  lv revoke-token lv_join_a3f9b2c1xxx

  Tip: run "lv invite --list" to see all active tokens`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tokenID := args[0]

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

		if cfg.VaultID == "" {
			return fmt.Errorf(
				"vault not registered on server\n  Run: lv push first",
			)
		}

		// Create signaling client
		sc := client.New(cfg.SignalingServer, id.DeviceID)
		if err := sc.HealthCheck(); err != nil {
			return err
		}

		// Revoke token
		if err := sc.RevokeToken(cfg.VaultID, tokenID); err != nil {
			return fmt.Errorf("failed to revoke token: %w", err)
		}

		fmt.Println()
		fmt.Println("╔══════════════════════════════════════════════════╗")
		fmt.Println("║           🚫 Token Revoked                       ║")
		fmt.Println("╠══════════════════════════════════════════════════╣")
		fmt.Printf("║  Token : %-39s║\n", tokenID[:20]+"...")
		fmt.Println("╠══════════════════════════════════════════════════╣")
		fmt.Println("║  ✅ Token is now invalid                         ║")
		fmt.Println("║  ✅ Nobody can join using this token anymore     ║")
		fmt.Println("║  ✅ Already joined peers keep their access       ║")
		fmt.Println("╠══════════════════════════════════════════════════╣")
		fmt.Println("║  To remove peer access use:                      ║")
		fmt.Println("║  lv revoke <device-id>                           ║")
		fmt.Println("╚══════════════════════════════════════════════════╝")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(revokeTokenCmd)
}
