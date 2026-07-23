package account

import "time"

// User is the subset of a "users" document the account domain reads/writes.
type User struct {
	ID               string    `bson:"_id"`
	Email            string    `bson:"email"`
	Name             string    `bson:"name"`
	AvatarURL        string    `bson:"avatar_url"`
	PasswordHash     string    `bson:"password_hash"`
	TwoFactorEnabled bool      `bson:"two_factor_enabled"`
	TwoFactorSecret  string    `bson:"two_factor_secret"`
	BackupCodes      []string  `bson:"backup_codes"` // sha256 hashes, single-use
	CreatedAt        time.Time `bson:"created_at"`
	UpdatedAt        time.Time `bson:"updated_at"`
}

// Session is the subset of a "sessions" document shown on the sessions page.
type Session struct {
	ID        string    `bson:"_id"`
	UserID    string    `bson:"user_id"`
	IP        string    `bson:"ip"`
	UserAgent string    `bson:"user_agent"`
	CreatedAt time.Time `bson:"created_at"`
	ExpiresAt time.Time `bson:"expires_at"`
}

// BackupCodeCount — how many one-time recovery codes we issue on 2FA verify.
const BackupCodeCount = 10
