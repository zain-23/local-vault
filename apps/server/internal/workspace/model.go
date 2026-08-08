package workspace

import "time"

type Workspace struct {
	ID        string    `bson:"_id" json:"id"`
	Name      string    `bson:"name" json:"name"`
	Slug      string    `bson:"slug" json:"slug"`
	Icon      string    `bson:"icon" json:"icon"`
	OwnerID   string    `bson:"owner_id" json:"owner_id"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// DefaultIcon is used when create/update omits icon, and for older docs with no icon.
const DefaultIcon = "vault"

// AllowedIcons is the allowlist of preset workspace icon keys.
var AllowedIcons = map[string]struct{}{
	"vault":       {},
	"lock":        {},
	"key":         {},
	"shield":      {},
	"folder":      {},
	"rocket":      {},
	"wrench":      {},
	"database":    {},
	"cloud":       {},
	"terminal":    {},
	"boxes":       {},
	"fingerprint": {},
}

// NormalizeIcon returns icon if allowed, otherwise DefaultIcon.
func NormalizeIcon(icon string) string {
	if _, ok := AllowedIcons[icon]; ok {
		return icon
	}
	return DefaultIcon
}

type Membership struct {
	ID          string    `bson:"_id" json:"id"`
	WorkspaceID string    `bson:"workspace_id" json:"workspace_id"`
	UserID      string    `bson:"user_id" json:"user_id"`
	Role        string    `bson:"role" json:"role"`    // owner | admin | member
	InvitedBy   string    `bson:"invited_by,omitempty" json:"invited_by,omitempty"`
	JoinedAt    time.Time `bson:"joined_at" json:"joined_at"`
}


const (
	RoleOwner	= "owner"
	RoleAdmin	= "admin"
	RoleMember	= "member"
)
