package id

import (
	"strings"
	"testing"
)

func TestGenerate_HasPrefix(t *testing.T) {
	got := Generate("usr_", 12)
	if !strings.HasPrefix(got, "usr_") {
		t.Errorf("expected prefix 'usr_', got %q", got) // %q prints with quotes
	}
}

func TestGenerate_CorrectLength(t *testing.T) {
	got := Generate("usr_", 12)
	// "usr_" = 4 chars + 12 bytes × 2 hex chars = 28 total
	if len(got) != 28 {
		t.Errorf("expected length 28, got %d (%q)", len(got), got)
	}
}

func TestGenerate_Unique(t *testing.T) {
	// two calls should never produce the same ID
	a := Generate("usr_", 12)
	b := Generate("usr_", 12)
	if a == b {
		t.Errorf("expected unique IDs, got same: %q", a)
	}
}