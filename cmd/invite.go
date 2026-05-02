package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/zain-23/local-vault/internal/client"
	"github.com/zain-23/local-vault/internal/config"
	"github.com/zain-23/local-vault/internal/identity"
	"github.com/zain-23/local-vault/internal/vault"
)

var inviteCmd = &cobra.Command{
	Use:   "invite",
	Short: "Generate an invite code to share with a teammate",
	RunE: func(cmd *cobra.Command, args []string) error {
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

		// Create signaling client
		sc := client.New(cfg.SignalingServer, id.DeviceID)

		fmt.Println("🔍 Checking signaling server...")
		if err := sc.HealthCheck(); err != nil {
			return err
		}

		// Generate random invite code
		code, err := generateInviteCode()
		if err != nil {
			return err
		}

		// Register invite on signaling server
		// Send X25519 public key so peer can encrypt secrets for us
		_, err = sc.CreateInvite(client.CreateInviteRequest{
			Code:      code,
			DeviceID:  id.DeviceID,
			PublicKey: id.X25519PublicKey, // X25519 for encryption
			IPHint:    getLocalIP(),
		})
		if err != nil {
			return err
		}

		// Show invite code to user
		fmt.Println()
		fmt.Println("🔐 LocalVault Invite")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Printf("  Invite Code: %s\n", code)
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println()
		fmt.Println("Share this code with your teammate.")
		fmt.Println("They should run:")
		fmt.Printf("  lv join %s\n", code)
		fmt.Println()
		fmt.Printf("⏳ Expires in 10 minutes (%s)\n",
			time.Now().Add(10*time.Minute).Format("15:04:05"))
		fmt.Printf("📱 Device: %s\n", id.DeviceName)
		fmt.Println()
		fmt.Println("⏳ Waiting for teammate to join...")
		fmt.Println("   (press Ctrl+C to cancel and run lv push manually)")
		fmt.Println()

		// Load vault so we can save peer when they join
		passphrase, err := promptPassphrase()
		if err != nil {
			return err
		}

		v, err := vault.Load(dir, passphrase)
		if err != nil {
			return err
		}

		// Poll signaling server every 3 seconds
		// Waiting for Dev B to join and send hello
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		// Stop waiting after 10 minutes
		timeout := time.After(10 * time.Minute)

		for {
			select {
			case <-timeout:
				fmt.Println()
				fmt.Println("⏰ Invite expired. Run lv invite again.")
				return nil

			case <-ticker.C:
				// Check mailbox for Dev B's hello
				msgs, err := sc.GetMessages()
				if err != nil {
					continue // network error — keep trying
				}

				if msgs.Count == 0 {
					fmt.Print(".") // show progress dots while waiting
					continue
				}

				// Someone joined
				fmt.Println()
				fmt.Println("🎉 Teammate joined!")

				for _, msg := range msgs.Messages {
					if msg.FromPublicKey == nil {
						continue
					}

					// Save them as trusted peer using their X25519 public key
					err = v.AddPeer(vault.Peer{
						DeviceID:        msg.FromDeviceID,
						DeviceName:      msg.FromDeviceID,
						PublicKey:       msg.FromPublicKey,
						X25519PublicKey: msg.FromPublicKey,
					})
					if err != nil {
						fmt.Printf("⚠️  Could not save peer: %v\n", err)
						continue
					}

					fmt.Printf("✅ Peer saved: %s\n", msg.FromDeviceID)
				}

				fmt.Println()
				fmt.Println("Now run: lv push")
				fmt.Println("To send your secrets to the new teammate.")
				return nil
			}
		}
	},
}

// generateInviteCode creates a random LV-XXXX-XXXX code
func generateInviteCode() (string, error) {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	encoded := base64.URLEncoding.EncodeToString(b)
	return fmt.Sprintf("LV-%s-%s", encoded[:4], encoded[4:8]), nil
}

// getLocalIP returns machine's local network IP
// Used as hint for LAN peer discovery
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok &&
			!ipnet.IP.IsLoopback() &&
			ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	return "127.0.0.1"
}

func init() {
	rootCmd.AddCommand(inviteCmd)
}
