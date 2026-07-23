package appstate

import (
	"testing"
)

func TestLoadSeedsStableFingerprint(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("SERVER_URL", "")

	first, err := Load()
	if err != nil {
		t.Fatalf("first load: %v", err)
	}
	if first.DeviceFingerprint == "" {
		t.Fatal("fingerprint not generated")
	}
	if first.DeviceName == "" {
		t.Fatal("device name not set")
	}

	second, err := Load()
	if err != nil {
		t.Fatalf("second load: %v", err)
	}
	if second.DeviceFingerprint != first.DeviceFingerprint {
		t.Fatalf("fingerprint changed: %q != %q", second.DeviceFingerprint, first.DeviceFingerprint)
	}
}

func TestServerURLEnvWins(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("SERVER_URL", "https://api.example.com")

	s, err := Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if s.ServerURL != "https://api.example.com" {
		t.Fatalf("want env server url, got %q", s.ServerURL)
	}
}

func TestServerURLEnvOverridesPersistedValue(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("SERVER_URL", "")

	first, err := Load()
	if err != nil {
		t.Fatalf("first load: %v", err)
	}
	if first.ServerURL != "http://localhost:8080" {
		t.Fatalf("want default server url, got %q", first.ServerURL)
	}

	t.Setenv("SERVER_URL", "https://api.example.com")
	second, err := Load()
	if err != nil {
		t.Fatalf("second load: %v", err)
	}
	if second.ServerURL != "https://api.example.com" {
		t.Fatalf("want env server url, got %q", second.ServerURL)
	}

	t.Setenv("SERVER_URL", "")
	third, err := Load()
	if err != nil {
		t.Fatalf("third load: %v", err)
	}
	if third.ServerURL != "https://api.example.com" {
		t.Fatalf("want persisted server url, got %q", third.ServerURL)
	}
}
