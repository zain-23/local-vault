package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/zain-23/local-vault/apps/cli/internal/ui"
)

var rootCmd = &cobra.Command{
	Use:   "lv",
	Short: "LocalVault — encrypted secrets for dev teams",
	Long:  "LocalVault\nEncrypted secret sync. No cloud. No leaks.",
}

// Execute is called from main.go
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		ui.Error("%s", err)
		os.Exit(1)
	}
}
