package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/zain-23/local-vault/apps/server/internal/common/apperror"
	"github.com/zain-23/local-vault/apps/server/internal/common/id"
	"github.com/zain-23/local-vault/apps/server/internal/common/jwt"
	"github.com/zain-23/local-vault/apps/server/internal/config"
)

// Service handles all auth business login
type Service struct {
	Store		*Store
	jwt 		*jwt.Service
	cfg			config.Config
}

func NewService(store *Store, jwtSvc *jwt.Service, cfg config.Config) *Service {
	return &Service{
		Store: store,
		jwt: jwtSvc,
		cfg: cfg,
	}
}

// generateRandomToken create URL-safe random string - used for email links
func generateRandomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	// URL-safe alphabet (-, _) so tokens survive email links / query strings —
	// standard base64's + and / get mangled (+ decodes to a space in a query).
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// sha25Hash hashes a token = we store the hash, never the raw token, so leaked DB can't reuse tokens
func sha256Hash(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:]) // h[:] converts fixed-size array to slice
}

// createSessionAndTokens generates tokens + store session
func (s *Service) createSessionAndTokens(ctx context.Context, user *User, deviceID, ip, userAgent string) (*LoginResponse, error) {
	sessionID := id.Generate("ses_", 12)

	accessToken, err := s.jwt.GenerateAccessToken(user.ID, user.Email, deviceID, sessionID)
	if err != nil {
		return  nil, apperror.ErrInternal
	}

	refreshToken, err := generateRandomToken()
	if err != nil {
		return nil, apperror.ErrInternal
	}

	now := time.Now()
	if err := s.Store.CreateSession(ctx, &Session{
		ID: sessionID,
		IP: ip,
		UserAgent: userAgent,
		UserID: user.ID,
		DeviceID: deviceID,
		RefreshTokenHash: sha256Hash(refreshToken),
		CreatedAt: now,
		ExpiresAt: now.Add(s.cfg.JWTRefreshExpiry),
	}); err != nil {
		return nil, apperror.ErrInternal
	}

	return &LoginResponse{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
		User: *user,
	}, nil
}

// IssueSession mints an access+refresh pair for an already-verified user.
// The device domain calls this after a user approves a CLI login — there is no
// password check here, because approval in an authenticated browser IS the proof.
// Exported wrapper so session creation lives in exactly one place.
func (s *Service) IssueSession(ctx context.Context, userID, deviceID, ip, userAgent string) (*LoginResponse, error) {
	user, err := s.Store.FindUserByID(ctx, userID)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	if user == nil {
		return nil, apperror.New(401, "user not found")
	}
	return s.createSessionAndTokens(ctx, user, deviceID, ip, userAgent)
}

// RevokeDeviceSessions deletes every session belonging to a device.
// Called when a device is removed, so revocation is immediate rather than
// waiting up to 15 minutes for the access token to expire.
func (s *Service) RevokeDeviceSessions(ctx context.Context, deviceID string) error {
	return s.Store.DeleteSessionsByDeviceID(ctx, deviceID)
}

func (s *Service) RefreshToken(ctx context.Context, req RefreshRequest) (*RefreshResponse, error) {
	session, err := s.Store.FindSessionByTokenHash(ctx, sha256Hash(req.RefreshToken))
	if err != nil {
		return nil, apperror.ErrInternal
	}
	if session == nil {
		return nil, apperror.New(401, "invalid or expired refresh token")
	}

	user, err := s.Store.FindUserByID(ctx, session.UserID)
	if err != nil || user == nil {
		return nil, apperror.New(401, "user not found")
	}

	accessToken, err := s.jwt.GenerateAccessToken(user.ID, user.Email, session.DeviceID, session.ID)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	return &RefreshResponse{AccessToken: accessToken}, nil
}


// ------------------------ Logout -----------------
func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	// Input already validated in the handler (RefreshRequest tag)
	return s.Store.DeleteSessionByTokenHash(ctx, sha256Hash(refreshToken))
}

// ---------------------- OAuth Helper -----------------
// FindOrCreateOAuthUser finds existing user or creates new one for OAuth login
func (s *Service) FindOrCreateOAuthUser(ctx context.Context, provider, oauthID, email, name, avatarURL string) (*User, error) {
	user, _ := s.Store.FindUserByOAuth(ctx, provider, oauthID)
	if user != nil {
		return user, nil
	}

	// Try by email - link OAuth to an existing account that hasn't linked this provider yet
	user, _ = s.Store.FindUserByEmail(ctx, strings.ToLower(email))

	if user != nil {
		if err := s.Store.UpdateUser(ctx, user.ID, bson.M{
			"oauth_provider": provider,
			"oauth_id": oauthID,
			"avatar_url": avatarURL,
			"updated_at": time.Now(),
		}); err != nil {
			return nil, apperror.ErrInternal
		}
		user.OAuthProvider = provider
		user.OAuthID = oauthID
		user.AvatarURL = avatarURL
		return user, nil
	}

	// New user - OAuth emails are pre-verified by the provider
	now := time.Now()
	user = &User{
		ID: id.Generate("usr", 12),
		Email: email,
		Name: name,
		OAuthProvider: provider,
		OAuthID: oauthID,
		AvatarURL: avatarURL,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.Store.CreateUser(ctx, user); err != nil {
		return nil, apperror.ErrInternal
	}

	return user, nil
}

// OAuthLogin creates session for an OAuth user — called after provider verification
func (s *Service) OAuthLogin(ctx context.Context, user *User, ip, userAgent string) (*LoginResponse, error) {
	return s.createSessionAndTokens(ctx, user, "", ip, userAgent)
}
