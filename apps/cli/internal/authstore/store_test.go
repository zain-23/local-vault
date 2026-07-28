package authstore

import (
	"errors"
	"testing"

	"github.com/zalando/go-keyring"
)

func TestSaveLoadClear(t *testing.T) {
	keyring.MockInit() // in-memory keychain; no OS interaction
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	if _, err := Load(); !errors.Is(err, ErrNoTokens) {
		t.Fatalf("want ErrNoTokens on empty store, got %v", err)
	}

	in := &Tokens{AccessToken: "acc", RefreshToken: "ref", DeviceID: "dev_1"}
	if err := Save(in); err != nil {
		t.Fatalf("save: %v", err)
	}

	got, err := Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if got.AccessToken != "acc" || got.RefreshToken != "ref" || got.DeviceID != "dev_1" {
		t.Fatalf("round-trip mismatch: %+v", got)
	}

	if err := Clear(); err != nil {
		t.Fatalf("clear: %v", err)
	}
	if _, err := Load(); !errors.Is(err, ErrNoTokens) {
		t.Fatalf("want ErrNoTokens after clear, got %v", err)
	}
}
