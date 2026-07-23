package account

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/zain-23/local-vault/server/internal/audit"
	"github.com/zain-23/local-vault/server/internal/common/apperror"
	"github.com/zain-23/local-vault/server/internal/common/password"
	"github.com/zain-23/local-vault/server/internal/common/totp"
)

type Service struct {
	store *Store
	audit audit.Recorder
}

func NewService(store *Store, recorder audit.Recorder) *Service {
	return &Service{store: store, audit: recorder}
}

// requireUser loads the caller or returns 404 — shared by every method.
func (s *Service) requireUser(ctx context.Context, userID string) (*User, error) {
	u, err := s.store.FindUserByID(ctx, userID)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	if u == nil {
		return nil, apperror.New(404, "user not found")
	}
	return u, nil
}

// Me returns the caller's profile.
func (s *Service) Me(ctx context.Context, userID string) (*ProfileResponse, error) {
	u, err := s.requireUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &ProfileResponse{
		ID: u.ID, Email: u.Email, Name: u.Name, AvatarURL: u.AvatarURL,
		TwoFactorEnabled: u.TwoFactorEnabled, CreatedAt: u.CreatedAt,
	}, nil
}

// UpdateProfile changes name and/or avatar. Only non-empty fields are applied.
func (s *Service) UpdateProfile(ctx context.Context, userID string, req UpdateProfileRequest) (*ProfileResponse, error) {
	fields := bson.M{"updated_at": time.Now()}
	if req.Name != "" {
		fields["name"] = req.Name
	}
	if req.AvatarURL != "" {
		fields["avatar_url"] = req.AvatarURL
	}
	if err := s.store.UpdateUser(ctx, userID, fields); err != nil {
		return nil, apperror.ErrInternal
	}

	s.audit.Record(ctx, audit.Entry{Action: "account.updated", TargetType: "user", TargetID: userID})
	return s.Me(ctx, userID)
}

// ChangePassword verifies the current password, sets a new one, and logs out every
// OTHER session (currentSID stays alive so the caller isn't kicked out).
func (s *Service) ChangePassword(ctx context.Context, userID, currentSID string, req ChangePasswordRequest) error {
	u, err := s.requireUser(ctx, userID)
	if err != nil {
		return err
	}
	if u.PasswordHash == "" {
		return apperror.New(400, "this account has no password (OAuth login)")
	}
	if !password.Verify(req.CurrentPassword, u.PasswordHash) {
		return apperror.New(401, "current password is incorrect")
	}

	newHash, err := password.Hash(req.NewPassword)
	if err != nil {
		return apperror.ErrInternal
	}
	if err := s.store.UpdateUser(ctx, userID, bson.M{"password_hash": newHash, "updated_at": time.Now()}); err != nil {
		return apperror.ErrInternal
	}

	// Security: a password change invalidates every other session.
	if err := s.store.DeleteSessionsExcept(ctx, userID, currentSID); err != nil {
		return apperror.ErrInternal
	}

	s.audit.Record(ctx, audit.Entry{Action: "account.password.changed", TargetType: "user", TargetID: userID})
	return nil
}


// Enable2FA generates a secret and returns it + the otpauth URL. 2FA is NOT active
// until Verify2FA succeeds — so a failed setup can't lock the user out.
func (s *Service) Enable2FA(ctx context.Context, userID string) (*Enable2FAResponse, error) {
	u, err := s.requireUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if u.TwoFactorEnabled {
		return nil, apperror.New(400, "2FA is already enabled")
	}

	secret, url, err := totp.GenerateSecret(u.Email)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	// Store the pending secret; enabled stays false until verify.
	if err := s.store.UpdateUser(ctx, userID, bson.M{"two_factor_secret": secret, "updated_at": time.Now()}); err != nil {
		return nil, apperror.ErrInternal
	}
	return &Enable2FAResponse{Secret: secret, OtpauthURL: url}, nil
}

// Verify2FA confirms the user can produce a valid code, activates 2FA, and returns
// one-time backup codes (shown once, stored hashed).
func (s *Service) Verify2FA(ctx context.Context, userID string, req Verify2FARequest) (*Verify2FAResponse, error) {
	u, err := s.requireUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if u.TwoFactorEnabled {
		return nil, apperror.New(400, "2FA is already enabled")
	}
	if u.TwoFactorSecret == "" {
		return nil, apperror.New(400, "start 2FA setup first")
	}
	if !totp.Validate(u.TwoFactorSecret, req.TOTPCode) {
		return nil, apperror.New(401, "invalid 2FA code")
	}

	plain, hashed, err := totp.GenerateBackupCodes(BackupCodeCount)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	if err := s.store.UpdateUser(ctx, userID, bson.M{
		"two_factor_enabled": true,
		"backup_codes":       hashed,
		"updated_at":         time.Now(),
	}); err != nil {
		return nil, apperror.ErrInternal
	}

	s.audit.Record(ctx, audit.Entry{Action: "auth.2fa.enabled", TargetType: "user", TargetID: userID})
	return &Verify2FAResponse{BackupCodes: plain}, nil // returned exactly once
}

// Disable2FA turns 2FA off after proving control via a TOTP code OR a backup code.
func (s *Service) Disable2FA(ctx context.Context, userID string, req Disable2FARequest) error {
	u, err := s.requireUser(ctx, userID)
	if err != nil {
		return err
	}
	if !u.TwoFactorEnabled {
		return apperror.New(400, "2FA is not enabled")
	}

	ok := false
	switch {
	case req.TOTPCode != "":
		ok = totp.Validate(u.TwoFactorSecret, req.TOTPCode)
	case req.BackupCode != "":
		// consume the backup code — matched means it existed and is now spent
		ok, err = s.store.ConsumeBackupCode(ctx, userID, totp.HashCode(req.BackupCode))
		if err != nil {
			return apperror.ErrInternal
		}
	default:
		return apperror.New(400, "a totp_code or backup_code is required")
	}
	if !ok {
		return apperror.New(401, "invalid code")
	}

	// Clear all 2FA state.
	if err := s.store.UpdateUser(ctx, userID, bson.M{
		"two_factor_enabled": false,
		"two_factor_secret":  "",
		"backup_codes":       []string{},
		"updated_at":         time.Now(),
	}); err != nil {
		return apperror.ErrInternal
	}

	s.audit.Record(ctx, audit.Entry{Action: "auth.2fa.disabled", TargetType: "user", TargetID: userID})
	return nil
}


// ListSessions returns the caller's active sessions, flagging the current one.
func (s *Service) ListSessions(ctx context.Context, userID, currentSID string) ([]SessionResponse, error) {
	sessions, err := s.store.ListSessions(ctx, userID)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	out := make([]SessionResponse, 0, len(sessions)) // len 0 → JSON [] not null
	for _, sess := range sessions {
		out = append(out, SessionResponse{
			ID: sess.ID, IP: sess.IP, UserAgent: sess.UserAgent,
			Current:   sess.ID == currentSID, // needs the sid claim from Task 2
			CreatedAt: sess.CreatedAt, ExpiresAt: sess.ExpiresAt,
		})
	}
	return out, nil
}

// RevokeSession deletes one of the caller's sessions (404 if it isn't theirs).
func (s *Service) RevokeSession(ctx context.Context, userID, sessionID string) error {
	sess, err := s.store.FindSessionByID(ctx, sessionID)
	if err != nil {
		return apperror.ErrInternal
	}
	if sess == nil || sess.UserID != userID { // don't confirm another user's session ids
		return apperror.New(404, "session not found")
	}
	if err := s.store.DeleteSession(ctx, sessionID); err != nil {
		return apperror.ErrInternal
	}
	return nil
}

// RevokeOtherSessions logs out everywhere except the caller's current session.
func (s *Service) RevokeOtherSessions(ctx context.Context, userID, currentSID string) error {
	if err := s.store.DeleteSessionsExcept(ctx, userID, currentSID); err != nil {
		return apperror.ErrInternal
	}
	s.audit.Record(ctx, audit.Entry{Action: "account.sessions.revoked", TargetType: "user", TargetID: userID})
	return nil
}
