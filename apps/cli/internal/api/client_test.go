package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/zain-23/local-vault/apps/cli/internal/authstore"
	"github.com/zalando/go-keyring"
)

// keyringMockInit isolates go-keyring for tests.
func keyringMockInit() { keyring.MockInit() }

// writeEnvelope mimics the server's success envelope.
func writeEnvelope(w http.ResponseWriter, status int, data any) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"data": data, "success": true, "message": "ok", "status": status,
	})
}

func TestDoUnwrapsEnvelope(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer acc" {
			t.Errorf("missing bearer, got %q", got)
		}
		writeEnvelope(w, 200, map[string]string{"email": "a@b.com"})
	}))
	defer srv.Close()

	c := New(srv.URL)
	c.UseTokens(&authstore.Tokens{AccessToken: "acc", RefreshToken: "ref"})

	var out struct {
		Email string `json:"email"`
	}
	if err := c.do(http.MethodGet, "/x", nil, &out, true); err != nil {
		t.Fatalf("do: %v", err)
	}
	if out.Email != "a@b.com" {
		t.Fatalf("want unwrapped data, got %+v", out)
	}
}

func TestDoRefreshesOn401(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/device/refresh":
			writeEnvelope(w, 200, map[string]string{"access_token": "acc2"})
		case "/x":
			n := atomic.AddInt32(&calls, 1)
			if n == 1 {
				w.WriteHeader(401)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": "expired"})
				return
			}
			if got := r.Header.Get("Authorization"); got != "Bearer acc2" {
				t.Errorf("retry used stale token: %q", got)
			}
			writeEnvelope(w, 200, map[string]string{"ok": "yes"})
		}
	}))
	defer srv.Close()

	keyring_MockNoop(t)
	c := New(srv.URL)
	c.UseTokens(&authstore.Tokens{AccessToken: "acc", RefreshToken: "ref"})

	var out struct {
		OK string `json:"ok"`
	}
	if err := c.do(http.MethodGet, "/x", nil, &out, true); err != nil {
		t.Fatalf("do: %v", err)
	}
	if out.OK != "yes" {
		t.Fatalf("want success after refresh, got %+v", out)
	}
}

func TestDoErrorBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "bad input"})
	}))
	defer srv.Close()

	c := New(srv.URL)
	err := c.do(http.MethodPost, "/x", map[string]string{"a": "b"}, nil, false)
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("want *APIError, got %T (%v)", err, err)
	}
	if apiErr.Status != 400 || apiErr.Message != "bad input" {
		t.Fatalf("unexpected APIError: %+v", apiErr)
	}
}

func keyring_MockNoop(t *testing.T) {
	t.Helper()
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	// go-keyring MockInit keeps Save in-memory; refresh persists there.
	keyringMockInit()
}
