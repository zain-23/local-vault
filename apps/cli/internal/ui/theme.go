package ui

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

// Website dark-palette tokens (frontend/src/styles.css).
var (
	accent = lipgloss.Color("#e5b567") // --primary (gold)
	green  = lipgloss.Color("#7bc894") // --success
	amber  = lipgloss.Color("#e5a93c") // --warning
	red    = lipgloss.Color("#e0715c") // --destructive
	muted  = lipgloss.Color("#8c8c8c") // --muted-foreground
)

// colorEnabled is computed once; tests may override it.
var colorEnabled = detectColor()

// detectColor reports whether ANSI styling should be emitted on stderr.
func detectColor() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	if os.Getenv("TERM") == "dumb" {
		return false
	}
	return term.IsTerminal(int(os.Stderr.Fd()))
}

// style returns a lipgloss style, or a no-op plain style when color is disabled.
func style(c lipgloss.Color, bold bool) lipgloss.Style {
	s := lipgloss.NewStyle()
	if !colorEnabled {
		return s
	}
	s = s.Foreground(c)
	if bold {
		s = s.Bold(true)
	}
	return s
}
