package member

import "time"

// Member = one user's role in one workspace ("memberships" collection)
type Membership struct {
	ID			string		`bson:"_id" json:"id"`
	WorkspaceID	string		`bson:"workspace_id" json:"workspace_id"`
	UserID		string		`bson:"user_id" json:"user_id"`
	Role		string		`bson:"role" json:"role"`
	InvitedBy	string		`bson:"invited_by,omitempty" json:"invited_by,omitempty"` // null for the owner
	JoinedAt	time.Time 	`bson:"joined_at" json:"joined_at"`
}

// Invite = a pending invitation to join a workspace ("workspace_invites" collection).
type Invite struct {
	ID          string    `bson:"_id" json:"id"`                 // prefix "inv_"
	WorkspaceID string    `bson:"workspace_id" json:"workspace_id"`
	Email       string    `bson:"email" json:"email"`            // lowercased recipient address
	Role        string    `bson:"role" json:"role"`              // admin | member
	InvitedBy   string    `bson:"invited_by" json:"invited_by"`  // user id who sent it
	TokenHash   string    `bson:"token_hash" json:"-"`           // sha256 of the raw token; never exposed
	Status      string    `bson:"status" json:"status"`          // pending | accepted
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	ExpiresAt   time.Time `bson:"expires_at" json:"expires_at"`  // TTL index deletes the doc after this
}

// UserInfo = the subset of a "users" document we need for display (member list + email).
type UserInfo struct {
	ID        string `bson:"_id"`
	Name      string `bson:"name"`
	Email     string `bson:"email"`
	AvatarURL string `bson:"avatar_url"`
}

// Role constants — the three access levels (avoid magic strings).
const (
	RoleOwner  = "owner"
	RoleAdmin  = "admin"
	RoleMember = "member"
)

// Invite status values.
const (
	StatusPending  = "pending"
	StatusAccepted = "accepted"
)
