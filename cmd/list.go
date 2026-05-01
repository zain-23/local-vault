package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/zain-23/local-vault/internal/vault"
)

// Terminal styles using lipgloss (like CSS)
var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7C3AED"))

	keyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981")) // green

	envBadgeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")). // gray
			Italic(true)
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all secrets (keys only, values hidden)",
	RunE: func(cmd *cobra.Command, args []string) error {
		passphrase, err := promptPassphrase()
		if err != nil {
			return err
		}

		dir, _ := os.Getwd()
		v, err := vault.Load(dir, passphrase)
		if err != nil {
			return err
		}

		secrets := v.List(envFlag)

		if len(secrets) == 0 {
			fmt.Println("No secrets found. Add one with: lv add KEY=VALUE")
			return nil
		}

		fmt.Println(headerStyle.Render("🔐 LocalVault Secrets"))
		fmt.Println()

		// Print each secret (key only, value masked)
		for _, s := range secrets {
			env := s.Env
			if env == "" {
				env = "all"
			}

			fmt.Printf(
				"  %s %s\n",
				keyStyle.Render(s.Key),
				envBadgeStyle.Render("["+env+"]"),
			)
		}

		fmt.Printf("\n%d secret(s) total\n", len(secrets))
		return nil
	},
}

func init() {
	listCmd.Flags().StringVarP(&envFlag, "env", "e", "", "filter by environment")
	rootCmd.AddCommand(listCmd)
}
