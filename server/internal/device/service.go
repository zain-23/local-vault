package device

import (
	"context"
	"time"

	"github.com/zain-23/local-vault/server/internal/auth"
	"github.com/zain-23/local-vault/server/internal/common/apperror"
	"github.com/zain-23/local-vault/server/internal/common/id"
	"github.com/zain-23/local-vault/server/internal/config"
)

// AuthIssuer is the slice of auth.Service this domain needs. Declaring the
// interface here (not in auth) keeps the dependency one-directional and lets
// tests pass a fake — the same trick middleware.MembershipChecker uses.
type AuthIssuer interface {
	IssueSession(ctx context.Context, userID, deviceID, ip, userAgent string) (*auth.LoginResponse, error)
	RevokeDeviceSessions(ctx context.Context, deviceID string) error
}

// Service holds the device domain's business logic.
type Service struct {
	Store   *Store
	auth    AuthIssuer
	cfg     config.Config
}

// NewService wires the device domain — called once at startup.
func NewService(store *Store, authSvc AuthIssuer, cfg config.Config) *Service {
	return &Service{Store: store, auth: authSvc, cfg: cfg}
}

// Authorize starts a login attempt. Unauthenticated — anyone can call it, which is
// why it returns nothing sensitive about the account and why the device_code it
// mints is the only thing that can later collect tokens.
func (s *Service) Authorize(ctx context.Context, req AuthorizeRequest, ip string) (*AuthorizeResponse, error) {
	userCode, err := generateUserCode()
	if err != nil {
		return nil, apperror.ErrInternal
	}
	deviceCode, err := generateDeviceCode()
	if err != nil {
		return nil, apperror.ErrInternal
	}

	now := time.Now()
	authReq := &AuthRequest{
		ID:             id.Generate("dar_", 12),
		UserCode:       userCode,
		DeviceCodeHash: sha256Hex(deviceCode), // only the hash is stored
		DeviceName:     req.DeviceName,
		Fingerprint:    req.DeviceFingerprint,
		IP:             ip,
		Status:         StatusPending,
		Consumed:       false,
		CreatedAt:      now,
		ExpiresAt:      now.Add(RequestTTL),
	}
	if err := s.Store.CreateAuthRequest(ctx, authReq); err != nil {
		return nil, apperror.ErrInternal
	}

	return &AuthorizeResponse{
		DeviceCode: deviceCode, // returned exactly once — never retrievable again
		UserCode:   userCode,
		VerifyURL:  s.cfg.FrontendURL + "/device",
		ExpiresIn:  int(RequestTTL.Seconds()),
		Interval:   PollInterval,
	}, nil
}

// ApprovalDetails powers the browser approval screen.
func (s *Service) ApprovalDetails(ctx context.Context, userCode string) (*ApprovalDetailsResponse, error) {
	req, err := s.Store.FindAuthRequestByUserCode(ctx, userCode)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	if req == nil {
		return nil, apperror.New(404, "device authorization request not found")
	}
	// The TTL index deletes expired docs eventually, not instantly — check the time
	// ourselves rather than trusting Mongo's background sweep to have run.
	if time.Now().After(req.ExpiresAt) {
		return nil, apperror.New(401, "device authorization request expired")
	}
	if req.Status != StatusPending {
		return nil, apperror.New(401, "device authorization request expired") // already approved or denied — nothing left to decide
	}

	return &ApprovalDetailsResponse{
		DeviceName: req.DeviceName,
		IP:         req.IP,
		Status:     req.Status,
		CreatedAt:  req.CreatedAt,
		ExpiresAt:  req.ExpiresAt,
	}, nil
}


// Decide approves or denies a pending request. userID comes from the browser's JWT —
// approving is an authenticated act, which is exactly what authorizes the CLI.
func (s *Service) Decide(ctx context.Context, userCode string, req DecisionRequest, userID string) error {
	authReq, err := s.Store.FindAuthRequestByUserCode(ctx, userCode)
	if err != nil {
		return apperror.ErrInternal
	}
	if authReq == nil {
		return apperror.New(404, "device authorization request not found")
	}
	if time.Now().After(authReq.ExpiresAt) {
		return apperror.New(401, "device authorization request expired")
	}
	if authReq.Status != StatusPending {
		return apperror.New(401, "device authorization request expired") // already decided — don't let a second tab flip it
	}

	// Deny is terminal and needs no workspace.
	if req.Action == ActionDeny {
		if err := s.Store.DenyAuthRequest(ctx, authReq.ID); err != nil {
			return apperror.ErrInternal
		}
		return nil
	}

      // --- approve ---
      // No workspace, no membership check: approving just links the device to the
      // account. Per-request RequireRole decides which workspace it may touch.
	now := time.Now()
	device := &Device{
		ID:           id.Generate("dev_", 12),
		UserID:       userID,
		Name:         authReq.DeviceName,
		Fingerprint:  authReq.Fingerprint,
		IP:           authReq.IP,
		LastSeenAt:   now,
		AuthorizedAt: now,
		CreatedAt:    now,
	}
	if err := s.Store.CreateDevice(ctx, device); err != nil {
		return apperror.ErrInternal
	}

	// Stamp the request. Its status:pending filter is the guard against a double
	// approve; if it lost the race we'd have an orphan device row, so clean up.
	if err := s.Store.ApproveAuthRequest(ctx, authReq.ID, userID, device.ID); err != nil {
		s.Store.DeleteDevice(ctx, device.ID)
		return apperror.New(401, "device authorization request expired")
	}

	// No tokens here — the browser must never hold the CLI's credentials.
	return nil
}

// Poll is how the CLI collects its tokens. It is unauthenticated: possession of the
// secret device_code IS the credential, which is why that code is 32 random bytes
// and never appears in a URL or a log line.
func (s *Service) Poll(ctx context.Context, req PollRequest, ip, userAgent string) (*PollResponse, error) {
	authReq, err := s.Store.FindAuthRequestByDeviceCodeHash(ctx, sha256Hex(req.DeviceCode))
	if err != nil {
		return nil, apperror.ErrInternal
	}
	
	if authReq == nil {
		return nil, apperror.New(404, "device authorization request not found")
	}

	if time.Now().After(authReq.ExpiresAt) {
		return nil, apperror.New(401, "device authorization request expired")
	}

	switch authReq.Status {
	case StatusPending:
		return &PollResponse{Status: StatusPending}, nil // CLI keeps looping
	case StatusDenied:
		return &PollResponse{Status: StatusDenied}, nil // CLI stops, prints "denied"
	}

	// --- approved ---
	// Claim the request before minting anything. Losing this race means another
	// poll already took the tokens, so this caller gets nothing.
	won, err := s.Store.ConsumeAuthRequest(ctx, authReq.ID)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	if !won {
		return nil, apperror.New(401, "device authorization request expired") // already collected — single-use by design
	}

	// Reuse the auth domain's session minting so refresh/logout behave identically
	// for CLI and web. deviceID lands in the JWT claim and on the session row.
	tokens, err := s.auth.IssueSession(ctx, authReq.UserID, authReq.DeviceID, ip, userAgent)
	if err != nil {
		return nil, err
	}

	return &PollResponse{
		Status:       StatusApproved,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		DeviceID:     authReq.DeviceID,
	}, nil
}

// ListDevices returns the caller's own devices in a workspace. Scoped to userID so
// one member cannot enumerate another's machines.
func (s *Service) ListDevices(ctx context.Context, userID string) ([]DeviceResponse, error) {
	devices, err := s.Store.ListDevices(ctx, userID)
	if err != nil {
		return nil, apperror.ErrInternal
	}

	// Map model → DTO so internal fields (fingerprint) never reach the API.
	out := make([]DeviceResponse, 0, len(devices)) // len 0, cap len(devices): JSON [] not null
	for _, d := range devices {
		out = append(out, DeviceResponse{
			ID:           d.ID,
			Name:         d.Name,
			IP:           d.IP,
			LastSeenAt:   d.LastSeenAt,
			AuthorizedAt: d.AuthorizedAt,
		})
	}
	return out, nil
}

// RevokeDevice removes a device and immediately kills its sessions.
func (s *Service) RevokeDevice(ctx context.Context, deviceID, userID string) error {
	device, err := s.Store.FindDeviceByID(ctx, deviceID)
	if err != nil {
		return apperror.ErrInternal
	}
	// Same 404 for "missing" and "not yours" — don't confirm another user's device IDs.
	if device == nil || device.UserID != userID {
		return apperror.New(404, "device not found")
	}

	// Sessions first: if this fails we stop, leaving a device row we can retry.
	// Deleting the device first and failing here would leave live tokens with no
	// device row to revoke through — unrevokable access.
	if err := s.auth.RevokeDeviceSessions(ctx, deviceID); err != nil {
		return apperror.ErrInternal
	}
	if err := s.Store.DeleteDevice(ctx, deviceID); err != nil {
		return apperror.ErrInternal
	}
	return nil
}