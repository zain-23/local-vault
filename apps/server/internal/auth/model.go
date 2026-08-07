package auth

import "time"

// User = a registered account in the "users" collection
type User struct {
	ID            string    `bson:"_id" json:"id"`                                  // MongoDB uses _id as primary key
	Email         string    `bson:"email" json:"email"`
	Name          string    `bson:"name" json:"name"`
	OAuthProvider string    `bson:"oauth_provider" json:"oauth_provider,omitempty"` // always "github" — omitempty skips empty values in JSON
	OAuthID       string    `bson:"oauth_id" json:"-"`                              // provider's user ID, hidden from API
	AvatarURL     string    `bson:"avatar_url" json:"avatar_url,omitempty"`
	Onboarded     bool      `bson:"onboarded" json:"onboarded"`
	CreatedAt     time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at" json:"updated_at"`
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
