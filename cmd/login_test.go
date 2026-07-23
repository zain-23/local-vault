package cmd

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/zain-23/local-vault/internal/api"
)

type fakePoller struct {
	responses []*api.DevicePollResponse
	errs      []error
	i         int
}

func (f *fakePoller) DevicePoll(string) (*api.DevicePollResponse, error) {
	r, e := f.responses[f.i], f.errs[f.i]
	f.i++
	return r, e
}

func TestPollForApprovalApproved(t *testing.T) {
	p := &fakePoller{
		responses: []*api.DevicePollResponse{
			{Status: "pending"},
			{Status: "approved", AccessToken: "a", RefreshToken: "r", DeviceID: "d"},
		},
		errs: []error{nil, nil},
	}
	now := func() time.Time { return time.Unix(0, 0) }
	tok, err := pollForApproval(context.Background(), p, "dc", time.Millisecond, time.Hour, now)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if tok.AccessToken != "a" || tok.DeviceID != "d" {
		t.Fatalf("unexpected tokens: %+v", tok)
	}
}

func TestPollForApprovalDenied(t *testing.T) {
	p := &fakePoller{
		responses: []*api.DevicePollResponse{{Status: "denied"}},
		errs:      []error{nil},
	}
	now := func() time.Time { return time.Unix(0, 0) }
	_, err := pollForApproval(context.Background(), p, "dc", time.Millisecond, time.Hour, now)
	if err == nil || err.Error() != "authorization denied" {
		t.Fatalf("want denied error, got %v", err)
	}
}

func TestPollForApprovalExpired(t *testing.T) {
	// Empty responses: if DevicePoll is called, the fake panics — proving the
	// deadline branch returns before any poll happens.
	p := &fakePoller{}
	calls := 0
	now := func() time.Time {
		calls++
		if calls == 1 {
			return time.Unix(0, 0) // deadline is computed from this call
		}
		return time.Unix(1000, 0) // every later call is well past the deadline
	}
	_, err := pollForApproval(context.Background(), p, "dc", time.Millisecond, time.Second, now)
	if err == nil || !errors.Is(err, errCodeExpired) {
		t.Fatalf("want expired error, got %v", err)
	}
}
