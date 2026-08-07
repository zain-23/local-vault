package account

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/zain-23/local-vault/apps/server/internal/audit"
	"github.com/zain-23/local-vault/apps/server/internal/common/apperror"
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
		Onboarded: u.Onboarded,
		CreatedAt: u.CreatedAt,
	}, nil
}

// UpdateProfile changes the name and/or onboarding flag. Only non-empty fields are applied.
func (s *Service) UpdateProfile(ctx context.Context, userID string, req UpdateProfileRequest) (*ProfileResponse, error) {
	fields := bson.M{"updated_at": time.Now()}
	if req.Name != "" {
		fields["name"] = req.Name
	}
	// Pointer, so this only writes when the client actually sent the field.
	if req.Onboarded != nil {
		fields["onboarded"] = *req.Onboarded
	}
	if err := s.store.UpdateUser(ctx, userID, fields); err != nil {
		return nil, apperror.ErrInternal
	}

	s.audit.Record(ctx, audit.Entry{Action: "account.updated", TargetType: "user", TargetID: userID})
	return s.Me(ctx, userID)
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
