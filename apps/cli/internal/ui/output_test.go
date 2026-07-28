package ui

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintersPlain(t *testing.T) {
	var buf bytes.Buffer
	out = &buf
	colorEnabled = false
	t.Cleanup(func() { out = defaultOut(); colorEnabled = detectColor() })

	cases := []struct {
		name string
		fn   func()
		want string
	}{
		{"success", func() { Success("done %d", 1) }, "✓ done 1\n"},
		{"error", func() { Error("nope") }, "✗ nope\n"},
		{"warn", func() { Warn("careful") }, "⚠ careful\n"},
		{"info", func() { Info("fyi") }, "› fyi\n"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buf.Reset()
			c.fn()
			if got := buf.String(); got != c.want {
				t.Fatalf("got %q, want %q", got, c.want)
			}
		})
	}
}

func TestHintIndents(t *testing.T) {
	var buf bytes.Buffer
	out = &buf
	colorEnabled = false
	t.Cleanup(func() { out = defaultOut(); colorEnabled = detectColor() })

	Hint("run: lv login")
	if got := buf.String(); !strings.Contains(got, "run: lv login") || !strings.HasPrefix(got, "  ") {
		t.Fatalf("hint not indented: %q", got)
	}
}
