package audit

import "time"

// Event = one recorded action ("audit_events" collection). Append-only.
type Event struct {
	ID          string         `bson:"_id" json:"id"`                                  // prefix "evt_"
	WorkspaceID string         `bson:"workspace_id" json:"workspace_id"`
	ActorID     string         `bson:"actor_id,omitempty" json:"actor_id,omitempty"`   // usr_ (empty on public join)
	Action      string         `bson:"action" json:"action"`                           // e.g. "vault.push"
	TargetType  string         `bson:"target_type,omitempty" json:"target_type,omitempty"`
	TargetID    string         `bson:"target_id,omitempty" json:"target_id,omitempty"`
	TargetName  string         `bson:"target_name,omitempty" json:"target_name,omitempty"`
	Details     map[string]any `bson:"details,omitempty" json:"details,omitempty"`     // free-form extra context
	DeviceID    string         `bson:"device_id,omitempty" json:"device_id,omitempty"` // dev_ from the JWT
	IP          string         `bson:"ip,omitempty" json:"ip,omitempty"`
	CreatedAt   time.Time      `bson:"created_at" json:"created_at"`
}

// Workspace roles (mirrors member.Role* — duplicated to avoid importing member).
const (
	RoleOwner  = "owner"
	RoleAdmin  = "admin"
	RoleMember = "member"
)
