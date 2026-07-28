package cmd

import (
	"os"

	"github.com/charmbracelet/lipgloss" // terminal styling
	"github.com/spf13/cobra"            // CLI framework

	"github.com/zain-23/local-vault/internal/ui"
)

// This is the style for our banner text
// lipgloss works like CSS for the terminal
var titleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#7C3AED")). // purple color
	MarginBottom(1)

// rootCmd is the base command
// When user just types "lv" with no subcommand, this runs
var rootCmd = &cobra.Command{
	Use:   "lv",
	Short: "LocalVault — encrypted secrets for dev teams",
	Long:  titleStyle.Render("🔐 LocalVault") + "\nEncrypted secret sync. No cloud. No leaks.",
}

// Execute is called from main.go
// It starts the entire CLI
// If any command fails, it prints the error and exits
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		ui.Error("%s", err)
		os.Exit(1)
	}
}
