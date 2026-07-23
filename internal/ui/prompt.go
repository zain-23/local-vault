package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

// Confirm asks a yes/no question on stderr and reads the answer from stdin.
func Confirm(question string) (bool, error) {
	fmt.Fprint(out, style(amber, true).Render(iconWarn)+" "+question+" (y/N): ")
	reader := bufio.NewReader(os.Stdin)
	resp, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	resp = strings.TrimSpace(strings.ToLower(resp))
	return resp == "y" || resp == "yes", nil
}

// Passphrase reads a passphrase without echo.
func Passphrase(label string) (string, error) {
	fmt.Fprint(out, label+": ")
	b, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Fprintln(out)
	if err != nil {
		return "", fmt.Errorf("failed to read passphrase: %w", err)
	}
	if len(b) == 0 {
		return "", fmt.Errorf("passphrase cannot be empty")
	}
	return string(b), nil
}
