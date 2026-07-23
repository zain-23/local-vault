package ui

import (
	"bytes"
	"strings"
	"testing"
)

func TestTablePlainShape(t *testing.T) {
	var buf bytes.Buffer
	out = &buf
	colorEnabled = false
	t.Cleanup(func() { out = defaultOut(); colorEnabled = detectColor() })

	Table([]string{"KEY", "ENV"}, [][]string{{"DATABASE_URL", "prod"}, {"API", "dev"}})
	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(lines) != 3 {
		t.Fatalf("want 3 lines, got %d: %q", len(lines), buf.String())
	}
	if !strings.HasPrefix(lines[0], "KEY") || !strings.Contains(lines[1], "DATABASE_URL") {
		t.Fatalf("unexpected table:\n%s", buf.String())
	}
}

func TestKeyValuePlain(t *testing.T) {
	var buf bytes.Buffer
	out = &buf
	colorEnabled = false
	t.Cleanup(func() { out = defaultOut(); colorEnabled = detectColor() })

	KeyValue("Email", "a@b.com")
	if got := buf.String(); !strings.Contains(got, "Email") || !strings.HasSuffix(got, ": a@b.com\n") {
		t.Fatalf("unexpected keyvalue: %q", got)
	}
}
