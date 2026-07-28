package device

import "time"

// AuthRequest = one in-flight CLI login attempt ("device_auth_requests" collection).
// Created unauthenticated by the CLI; filled in when a signed-in user approves it.
type AuthRequest struct {
	ID             string    `bson:"_id" json:"id"`                          // prefix "dar_"
	UserCode       string    `bson:"user_code" json:"user_code"`             // "ABCD-1234" — shown to the human
	DeviceCodeHash string    `bson:"device_code_hash" json:"-"`              // sha256 of the CLI's secret; json:"-" = never exposed
	DeviceName     string    `bson:"device_name" json:"device_name"`
	Fingerprint    string    `bson:"device_fingerprint" json:"device_fingerprint"`
	IP             string    `bson:"ip" json:"ip"`
	UserID         string    `bson:"user_id,omitempty" json:"user_id,omitempty"`             // null until approved
	DeviceID       string    `bson:"device_id,omitempty" json:"device_id,omitempty"`         // null until approved
	Status         string    `bson:"status" json:"status"`                   // pending | approved | denied
	Consumed       bool      `bson:"consumed" json:"consumed"`               // true once the CLI has collected its tokens
	CreatedAt      time.Time `bson:"created_at" json:"created_at"`
	ExpiresAt      time.Time `bson:"expires_at" json:"expires_at"`           // TTL index deletes the doc after this
}

// Device = an authorized CLI installation ("devices" collection).
type Device struct {
	ID           string    `bson:"_id" json:"id"`                    // prefix "dev_"
	UserID       string    `bson:"user_id" json:"user_id"`
	Name         string    `bson:"name" json:"name"`
	Fingerprint  string    `bson:"fingerprint" json:"fingerprint"`
	IP           string    `bson:"ip" json:"ip"`
	LastSeenAt   time.Time `bson:"last_seen_at" json:"last_seen_at"`
	AuthorizedAt time.Time `bson:"authorized_at" json:"authorized_at"`
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
}

// Auth request status values.
const (
	StatusPending  = "pending"
	StatusApproved = "approved"
	StatusDenied   = "denied"
)

// Approve/deny actions accepted by PUT /device/authorize/:userCode.
const (
	ActionApprove = "approve"
	ActionDeny    = "deny"
)

// Flow timings — TTL for a request, and how often we tell the CLI to poll.
const (
	RequestTTL   = 10 * time.Minute
	PollInterval = 2 // seconds
)