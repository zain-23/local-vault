package ui

import (
	"fmt"
	"strings"
)

// Divider prints a themed horizontal rule.
func Divider() {
	fmt.Fprintln(out, style(muted, false).Render(strings.Repeat("─", 40)))
}

// Header prints a section heading followed by a divider.
func Header(text string) {
	fmt.Fprintln(out, style(accent, true).Render(text))
	Divider()
}

// KeyValue prints an aligned "  key : value" row.
func KeyValue(key, value string) {
	label := style(muted, false).Render(fmt.Sprintf("%-12s", key))
	fmt.Fprintln(out, "  "+label+": "+value)
}

// Table prints headers + rows with columns padded to the widest cell.
func Table(headers []string, rows [][]string) {
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, r := range rows {
		for i, cell := range r {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}
	line := func(cells []string, header bool) {
		var b strings.Builder
		for i, cell := range cells {
			w := 0
			if i < len(widths) {
				w = widths[i]
			}
			b.WriteString(fmt.Sprintf("%-*s  ", w, cell))
		}
		text := strings.TrimRight(b.String(), " ")
		if header {
			text = style(muted, true).Render(text)
		}
		fmt.Fprintln(out, text)
	}
	line(headers, true)
	for _, r := range rows {
		line(r, false)
	}
}
