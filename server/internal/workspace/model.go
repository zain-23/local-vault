package workspace

import "time"

type Workspace struct {
	ID			string 		`bson:"_id" json:"id"`
	Name		string		`bson:"name" json:"name"`
	Slug		string		`bson:"slug" json:"slug"`
	OwnerID		string		`bson:"owner_id" json:"owner_id"`
	CreatedAt	time.Time	`bson:"created_at" json:"created_at"`
	UpdatedAt 	time.Time	`bson:"updated_at" json:"updated_at"`
}

type Membership struct {
	ID          string    `bson:"_id" json:"id"`
	WorkspaceID string    `bson:"workspace_id" json:"workspace_id"`
	UserID      string    `bson:"user_id" json:"user_id"`
	Role        string    `bson:"role" json:"role"`    // owner | admin | member
	InvitedBy   string    `bson:"invited_by" json:"invited_by,omitempty"`
	JoinedAt    time.Time `bson:"joined_at" json:"joined_at"`
}


const (
	RoleOwner	= "owner"
	RoleAdmin	= "admin"
	RoleMember	= "member"
)
