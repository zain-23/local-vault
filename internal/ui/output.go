package ui

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/lipgloss"
)

// out is where all human-facing messaging is written. Default: stderr.
var out io.Writer = defaultOut()

func defaultOut() io.Writer { return os.Stderr }

func printLine(c lipgloss.Color, icon, format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Fprintln(out, style(c, true).Render(icon)+" "+msg)
}

// Success/Error/Warn/Info/Step print a single glyph-prefixed status line.
func Success(format string, a ...any) { printLine(green, iconSuccess, format, a...) }
func Error(format string, a ...any)   { printLine(red, iconError, format, a...) }
func Warn(format string, a ...any)    { printLine(amber, iconWarn, format, a...) }
func Info(format string, a ...any)    { printLine(muted, iconInfo, format, a...) }
func Step(format string, a ...any)    { printLine(muted, iconInfo, format, a...) }

// Hint prints a dim, indented follow-up suggestion.
func Hint(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Fprintln(out, style(muted, false).Render("  "+msg))
}

// Title prints the brand wordmark followed by a heading.
func Title(text string) {
	fmt.Fprintln(out, style(accent, true).Render(wordmark)+"  "+text)
}

// Code prints a value prominently (used for the device user-code).
func Code(text string) {
	fmt.Fprintln(out, "  "+style(accent, true).Render(text))
}
