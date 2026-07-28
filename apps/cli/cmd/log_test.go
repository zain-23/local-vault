package cmd

import "testing"

func TestShortID(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"local sentinel is shorter than 8 chars", "local", "local"},
		{"uuid is truncated to 8 chars plus ellipsis", "550e8400-e29b-41d4-a716-446655440000", "550e8400..."},
		{"exactly 8 chars is returned as-is", "abcdefgh", "abcdefgh"},
		{"empty string is returned as-is", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shortID(tt.in); got != tt.want {
				t.Errorf("shortID(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
