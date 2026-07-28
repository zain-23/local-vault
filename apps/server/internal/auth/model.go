package auth

import "time"

// User = a registered account in the "users" collection
type User struct {
	ID               string    `bson:"_id" json:"id"`                                  // MongoDB uses _id as primary key
	Email            string    `bson:"email" json:"email"`
	Name             string    `bson:"name" json:"name"`
	PasswordHash     string    `bson:"password_hash" json:"-"`                         // json:"-" = NEVER include in API responses (security)
	OAuthProvider    string    `bson:"oauth_provider" json:"oauth_provider,omitempty"`  // "google"/"github"/"" — omitempty skips empty values in JSON
	OAuthID          string    `bson:"oauth_id" json:"-"`                              // provider's user ID, hidden from API
	AvatarURL        string    `bson:"avatar_url" json:"avatar_url,omitempty"`
	EmailVerified    bool      `bson:"email_verified" json:"email_verified"`
	TwoFactorEnabled bool      `bson:"two_factor_enabled" json:"two_factor_enabled"`
	Onboarded		 bool	   `bson:"onboarded" json:"onboarded"`
	TwoFactorSecret  string    `bson:"two_factor_secret" json:"-"`                     // TOTP secret, never expose
	BackupCodes      []string  `bson:"backup_codes" json:"-"`                          // hashed backup codes
	CreatedAt        time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time `bson:"updated_at" json:"updated_at"`
}

// Session = one login session, stored in "sessions" collection
type Session struct {
	ID               string    `bson:"_id"`
	UserID           string    `bson:"user_id"`
	DeviceID		 string    `bson:"device_id,omitempty"`
	RefreshTokenHash string    `bson:"refresh_token_hash"` // hash of refresh token, not the token itself
	IP               string    `bson:"ip"`
	UserAgent        string    `bson:"user_agent"`
	CreatedAt        time.Time `bson:"created_at"`
	ExpiresAt        time.Time `bson:"expires_at"`         // MongoDB TTL index auto-deletes when this passes
}

// EmailVerification = token hash for email verification link
type EmailVerification struct {
	ID        string    `bson:"_id"`
	UserID    string    `bson:"user_id"`
	TokenHash string    `bson:"token_hash"`
	CreatedAt time.Time `bson:"created_at"`
	ExpiresAt time.Time `bson:"expires_at"`
}

// PasswordReset = token hash for password reset link
type PasswordReset struct {
	ID        string    `bson:"_id"`
	UserID    string    `bson:"user_id"`
	TokenHash string    `bson:"token_hash"`
	Used      bool      `bson:"used"` 
	CreatedAt time.Time `bson:"created_at"`
	ExpiresAt time.Time `bson:"expires_at"`
}

// MagicLink = token hash for passwordless login link
type MagicLink struct {
	ID        string    `bson:"_id"`
	Email     string    `bson:"email"`
	TokenHash string    `bson:"token_hash"`
	Used      bool      `bson:"used"`
	CreatedAt time.Time `bson:"created_at"`
	ExpiresAt time.Time `bson:"expires_at"`
}