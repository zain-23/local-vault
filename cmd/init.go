package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/zain-23/local-vault/internal/client"
	"github.com/zain-23/local-vault/internal/config"
	"github.com/zain-23/local-vault/internal/identity"
	"github.com/zain-23/local-vault/internal/session"
	"github.com/zain-23/local-vault/internal/vault"
	"golang.org/x/term"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new vault in the current directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔐 Initializing LocalVault...")

		// Ask passphrase
		fmt.Print("Enter passphrase: ")
		passphraseBytes, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			return fmt.Errorf("failed to read passphrase: %w", err)
		}

		fmt.Print("Confirm passphrase: ")
		confirmBytes, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			return fmt.Errorf("failed to read passphrase: %w", err)
		}

		if string(passphraseBytes) != string(confirmBytes) {
			return fmt.Errorf("passphrases do not match")
		}

		passphrase := string(passphraseBytes)
		if len(passphrase) < 8 {
			return fmt.Errorf("passphrase must be at least 8 characters")
		}

		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		lvDir := filepath.Join(dir, ".lv")

		// Initialize vault on disk
		if err := vault.Init(dir, passphrase); err != nil {
			return err
		}

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

		// Register vault on signaling server
		fmt.Println("🌐 Registering vault on server...")
		sc := client.New(cfg.SignalingServer, id.DeviceID)

		if err := sc.HealthCheck(); err != nil {
			fmt.Println("⚠️  Could not reach server — working offline")
			fmt.Println("   Run: lv push when online to register")
		} else {
			resp, err := sc.RegisterVault(client.RegisterVaultRequest{
				OwnerID:         id.DeviceID,
				OwnerName:       id.DeviceName,
				PublicKey:       id.PublicKey,
				X25519PublicKey: id.X25519PublicKey,
			})

			if err != nil {
				fmt.Printf("⚠️  Server registration failed: %v\n", err)
			} else {
				// Save vault ID to config
				cfg.VaultID = resp.VaultID
				cfg.DeviceID = id.DeviceID
				config.Save(lvDir, cfg)

				fmt.Printf("✅ Registered — Vault ID: %s\n", resp.VaultID)
			}
		}

		// Auto unlock after init
		v, err := vault.Load(dir, passphrase)
		if err == nil {
			if err := session.Save(lvDir, v.GetKey()); err == nil {
				fmt.Println("🔓 Auto-unlocked for 12 hours")
			}
		}

		fmt.Println("✅ Vault initialized successfully")
		fmt.Println("📁 Created .lv/ directory")
		fmt.Println("🔒 Vault encrypted with your passphrase")
		fmt.Println("📝 Updated .gitignore")
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Println("  lv add DATABASE_URL=postgres://...")
		fmt.Println("  lv push")
		fmt.Println("  lv invite --name \"Ahmed\"")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
