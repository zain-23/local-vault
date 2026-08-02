package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/zain-23/local-vault/apps/server/internal/common/apperror"
	"github.com/zain-23/local-vault/apps/server/internal/common/id"
	"github.com/zain-23/local-vault/apps/server/internal/common/jwt"
	"github.com/zain-23/local-vault/apps/server/internal/common/password"
	"github.com/zain-23/local-vault/apps/server/internal/common/totp"
	"github.com/zain-23/local-vault/apps/server/internal/config"
	"github.com/zain-23/local-vault/apps/server/internal/email"
)

// Argon2id params - controls how hard it is to brute-force password
const (
	argonTime		= 2				// iteration - more = slower but safer
	argonMemory		= 64 * 1024		// 64MB RAM - makes GPU attacks expensive
	argonThreads 	= 4				// parallel threads
	argonKeyLen		= 32			// output hash size in bytes
	argonSaltLen	= 16			// random salt size
)

// Service handles all auth business login
type Service struct {
	Store		*Store
	jwt 		*jwt.Service
	cfg			config.Config
	publisher	*email.Publisher
}

func NewService(store *Store, jwtSvc *jwt.Service, pub *email.Publisher, cfg config.Config) *Service {
	return &Service{
		Store: store,
		jwt: jwtSvc,
		cfg: cfg,
		publisher: pub,
	}
}


// HashPassword converts plain password to safe hash
func hashPassword(pw string) (string, error) {
	return password.Hash(pw)
}

// verifyPassword re-hashes with same salth and compares - returns true if password matches
func verifyPassword(pw, encode string) bool {
	return password.Verify(pw, encode)
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

// verifyTOTP — placeholder until Account domain adds 2FA setup
func verifyTOTP(secret, code string) bool {
	return false // no user has 2FA enabled yet — setup endpoints come in Account domain
}

// ------------- Signup ------------------------
func (s *Service) Signup(ctx context.Context, req SignupRequest) (string, error) {
	// reject duplicate emails
	existing, err := s.Store.FindUserByEmail(ctx, strings.ToLower(req.Email))
	if err != nil {
		return "", apperror.ErrInternal
	}

	if existing != nil {
		return  "", apperror.New(409, "email already registered")
	}

	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return  "", apperror.ErrInternal
	}

	now := time.Now()
	user := &User{
		ID: 			id.Generate("usr_", 12),
		Email: 			strings.ToLower(req.Email),
		Name: 			req.Name,
		PasswordHash: 	passwordHash,
		EmailVerified: 	false,
		Onboarded: 		false,
		CreatedAt: 		now,
		UpdatedAt: 		now,	
	}

	if err := s.Store.CreateUser(ctx, user); err != nil {
		return "", apperror.ErrInternal
	}

	// send verification email - user must prove they own this email
	token, err := generateRandomToken()
	if err != nil {
		return  "", apperror.ErrInternal
	}
	s.Store.CreateEmailVerification(ctx, &EmailVerification{
		ID: 		id.Generate("evr_", 12),
		UserID: 	user.ID,
		TokenHash: 	sha256Hash(token),
		CreatedAt: 	now,
		ExpiresAt: 	now.Add(24 * time.Hour),
	})

	// Build the link and enqueue the email - a queue hiccup must not fail signup
	verifyURL := fmt.Sprintf("%s/auth/verify-email?token=%s", s.cfg.FrontendURL, token)
	job := email.EmailJob{
		Kind: email.KindVerification,
		Name: user.Name,
		URL: verifyURL,
		To: user.Email,
	}

	if err := s.publisher.Publish(ctx, job); err != nil {
		log.Printf("⚠️ failed to enqueue verification email for %s: %v", user.Email, err)
	}
	return  "Account created. Check your email to verify.", nil
}

// -------------------- Login ---------------------
func (s *Service) Login(ctx context.Context, req LoginRequest, ip, userAgent string) (*LoginResult, error) {
	user, err := s.Store.FindUserByEmail(ctx, strings.ToLower(req.Email))
	if err != nil {
		return  nil, apperror.ErrInternal
	}

	if user == nil {
		return nil, apperror.New(401, "invalid email or password")
	}

	// OAuth-only accounts have no password - they must login via OAuth
	if user.PasswordHash == "" {
		return nil, apperror.New(401, "this account uses OAuth login")
	}

	if !verifyPassword(req.Password, user.PasswordHash) {
		return  nil, apperror.New(401, "invalid email or password")
	}

	if !user.EmailVerified {
		return  nil, apperror.New(403, "please verify your email before logging in")
	}

	// if 2FA enabled - don't give full tokens yet, require TOTP code first
	if user.TwoFactorEnabled {
		tempToken, err := s.jwt.GenerateTempToken(user.ID)
		if err != nil {
			return  nil, apperror.ErrInternal
		}

		return &LoginResult{
			Requires2FA: true,
			TempToken: tempToken,
		}, nil
	}

	tokens, err := s.createSessionAndTokens(ctx, user, "", ip, userAgent)
	if err != nil {
		return nil, err
	}

	return  &LoginResult{Tokens: tokens}, nil
} 

// Login2FA completes login after the user proves 2FA with a TOTP or backup code.
func (s *Service) Login2FA(ctx context.Context, req Login2FARequest, ip, userAgent string) (*LoginResponse, error) {
	claims, err := s.jwt.ValidateToken(req.TempToken)
	if err != nil || claims.Type != "2fa_temp" {
		return nil, apperror.New(401, "invalid or expired 2FA token")
	}

	user, err := s.Store.FindUserByID(ctx, claims.Subject)
	if err != nil || user == nil {
		return nil, apperror.New(401, "user not found")
	}
	if !user.TwoFactorEnabled {
		return nil, apperror.New(400, "2FA is not enabled")
	}

	ok := false
	switch {
	case req.TOTPCode != "":
		ok = totp.Validate(user.TwoFactorSecret, req.TOTPCode)
	case req.BackupCode != "":
		ok, err = s.Store.ConsumeBackupCode(ctx, user.ID, totp.HashCode(req.BackupCode))
		if err != nil {
			return nil, apperror.ErrInternal
		}
	default:
		return nil, apperror.New(400, "a totp_code or backup_code is required")
	}
	if !ok {
		return nil, apperror.New(401, "invalid 2FA code")
	}

	return s.createSessionAndTokens(ctx, user, "", ip, userAgent)
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

// ---------------------- Email verification ------------------
func (s *Service) VerifyEmail(ctx context.Context, token string) error {
	ev, err := s.Store.FindEmailVerificationByHash(ctx, sha256Hash(token))
	if err != nil {
		return apperror.ErrInternal
	}
	if ev == nil {
		return apperror.New(400, "invalid or expired verification token")
	}

	// Mark verified only — no auto-login, the user logs in themselves
	s.Store.UpdateUser(ctx, ev.UserID, bson.M{"email_verified": true, "updated_at": time.Now()})
	s.Store.DeleteEmailVerification(ctx, ev.ID)
	return nil
}

// ForgotPassword always returns success — prevents attackers from discovering which emails exist
func (s *Service) ForgotPassword(ctx context.Context, req ForgotPasswordRequest) (string, error) {
	user, _ := s.Store.FindUserByEmail(ctx, strings.ToLower(req.Email))
	if user != nil {
		token, _ := generateRandomToken()
		now := time.Now()
		s.Store.CreatePasswordReset(ctx, &PasswordReset{
			ID: id.Generate("prs_", 12), UserID: user.ID,
			TokenHash: sha256Hash(token), Used: false,
			CreatedAt: now, 
			ExpiresAt: now.Add(1 * time.Hour),
		})
		resetURL := fmt.Sprintf("%s/auth/reset-password?token=%s", s.cfg.FrontendURL, token)
		job := email.EmailJob{
			Kind: email.KindPasswordReset,
			To: user.Email,
			Name: user.Name,
			URL: resetURL,
		}

		if err := s.publisher.Publish(ctx, job); err != nil {
			log.Printf("⚠️ failed to enqueue password forgot email for %s: %v", user.Email, err)
		}
	}
	return "If that email exists, a reset link has been sent.", nil
}

func (s *Service) ResetPassword(ctx context.Context, req ResetPasswordRequest) (string, error) {
	// Input already validated in the handler (token required, new_password min 8)
	pr, err := s.Store.FindPasswordResetByHash(ctx, sha256Hash(req.Token))
	if err != nil {
		return "", apperror.ErrInternal
	}
	if pr == nil {
		return "", apperror.New(400, "invalid or expired reset token")
	}

	passwordHash, err := hashPassword(req.NewPassword)
	if err != nil {
		return "", apperror.ErrInternal
	}
	s.Store.UpdateUser(ctx, pr.UserID, bson.M{"password_hash": passwordHash, "updated_at": time.Now()})
	s.Store.MarkPasswordResetUsed(ctx, pr.ID)

	return "Password reset successfully. You can now log in.", nil
}

// ---------------------------- Magic Link --------------------
// SendMagicLink always returns success - same reason as ForgotPassword
func (s *Service) SendMagicLink(ctx context.Context, req MagicLinkRequest) (string, error) {
	user, _ := s.Store.FindUserByEmail(ctx, strings.ToLower(req.Email))
	
	if user != nil {
		token, _ := generateRandomToken()
		now := time.Now()
		s.Store.CreateMagicLink(ctx, &MagicLink{
			ID: id.Generate("mgl_", 12),
			Email: user.Email,
			TokenHash: sha256Hash(token),
			Used: false,
			CreatedAt: now,
			ExpiresAt: now.Add(15 * time.Minute),
		})

		magicURL := fmt.Sprintf("%s/magic-link?token=%s", s.cfg.FrontendURL, token)
		job := email.EmailJob{
			Kind: email.KindSendMagicLink,
			To: user.Email,
			Name: user.Name,
			URL: magicURL,
		}
		if err := s.publisher.Publish(ctx, job); err != nil {
			log.Printf("⚠️ failed to enqueue password send magic email for %s: %v", user.Email, err)
		}
	}
	return "If that email exists, a login link has been sent.", nil
}

func (s *Service) VerifyMagicLink(ctx context.Context, req MagicLinkVerifyRequest, ip, userAgent string) (*LoginResponse, error) {
	// Input already validated in the handler
	ml, err := s.Store.FindMagicLinkByHash(ctx, sha256Hash(req.Token))
	if err != nil {
		return nil, apperror.ErrInternal
	}
	if ml == nil {
		return nil, apperror.New(400, "invalid or expired magic link")
	}

	s.Store.MarkMagicLinkUsed(ctx, ml.ID)

	user, err := s.Store.FindUserByEmail(ctx, ml.Email)
	if err != nil || user == nil {
		return nil, apperror.ErrInternal
	}
	// Auto-verify email via magic link — user proved they have access
	if !user.EmailVerified {
		s.Store.UpdateUser(ctx, user.ID, bson.M{"email_verified": true, "updated_at": time.Now()})
	}
	return s.createSessionAndTokens(ctx, user, "", ip, userAgent)
}


// ---------------------- OAuth Helper -----------------
// FindOrCreateOAuthUser finds existing user or creates new one for OAuth login
func (s *Service) FindOrCreateOAuthUser(ctx context.Context, provider, oauthID, email, name, avatarURL string) (*User, error) {
	user, _ := s.Store.FindUserByOAuth(ctx, provider, oauthID)
	if user != nil {
		return user, nil
	}

	// Try by email - link OAuth to existing email/password account
	user, _ = s.Store.FindUserByEmail(ctx, strings.ToLower(email))

	if user != nil {
		s.Store.UpdateUser(ctx, user.ID, bson.M{
			"oauth_provider": provider,
			"oauth_id": oauthID,
			"avatar_url": avatarURL,
			"email_verified": true,
			"updated_at": time.Now(),
		})
		return user, bson.ErrDecodeToNil
	}

	// New user - OAuth emails are pre-verified by the provider
	now := time.Now()
	user = &User{
		ID: id.Generate("usr", 12),
		Email: email,
		Name: name,
		OAuthProvider: provider,
		OAuthID: oauthID,
		EmailVerified: true,
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
