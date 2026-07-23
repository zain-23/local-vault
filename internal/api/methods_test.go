package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zain-23/local-vault/internal/authstore"
)

func TestDeviceAuthorizeAndPoll(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/device/authorize":
			if r.Header.Get("Authorization") != "" {
				t.Error("authorize must be unauthenticated")
			}
			writeEnvelope(w, 201, map[string]any{
				"device_code": "dc", "user_code": "WXYZ-1234",
				"verify_url": "http://f/device", "expires_in": 600, "interval": 5,
			})
		case "/api/v1/device/authorize/poll":
			writeEnvelope(w, 200, map[string]any{
				"status": "approved", "access_token": "a", "refresh_token": "r", "device_id": "dev_1",
			})
		}
	}))
	defer srv.Close()

	c := New(srv.URL)
	auth, err := c.DeviceAuthorize("laptop", "fp-123")
	if err != nil {
		t.Fatalf("authorize: %v", err)
	}
	if auth.UserCode != "WXYZ-1234" || auth.Interval != 5 {
		t.Fatalf("unexpected authorize response: %+v", auth)
	}

	poll, err := c.DevicePoll("dc")
	if err != nil {
		t.Fatalf("poll: %v", err)
	}
	if poll.Status != "approved" || poll.AccessToken != "a" {
		t.Fatalf("unexpected poll response: %+v", poll)
	}
}

func TestMe(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/account/me" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		writeEnvelope(w, 200, map[string]any{
			"id": "u1", "email": "a@b.com", "name": "Ada", "two_factor_enabled": true,
		})
	}))
	defer srv.Close()

	c := New(srv.URL)
	tokens := mustTokens()
	c.UseTokens(&tokens)
	acct, err := c.Me()
	if err != nil {
		t.Fatalf("me: %v", err)
	}
	if acct.Name != "Ada" || !acct.TwoFactorEnabled {
		t.Fatalf("unexpected account: %+v", acct)
	}
}

func mustTokens() authstore.Tokens {
	return authstore.Tokens{AccessToken: "a", RefreshToken: "r", DeviceID: "d"}
}
