package audit

import (
	"time"

	"github.com/zain-23/local-vault/apps/server/internal/common/pagination"
)

// Entry is what a domain service hands the recorder. Actor/ip/device are filled
// by the recorder from the request context — the domain only knows action + target.
type Entry struct {
	WorkspaceID string
	Action      string
	TargetType  string
	TargetID    string
	TargetName  string
	Details     map[string]any
}

// ListQuery is the GET /audit query string. Pagination is embedded; every filter
// is optional (omitempty → absent means "no filter").
type ListQuery struct {
	pagination.Params
	Action       string `query:"action"`
	ActionPrefix string `query:"action_prefix"` // e.g. "vault" → all vault.* events
	ActorID      string `query:"actor_id"`
	DeviceID     string `query:"device_id"`
	From         string `query:"from"` // RFC3339 lower bound (inclusive)
	To           string `query:"to"`   // RFC3339 upper bound (inclusive)
}

// Filter is the parsed, validated form of ListQuery used by the store.
type Filter struct {
	Action       string
	ActionPrefix string
	ActorID      string
	DeviceID     string
	From         *time.Time // nil = unbounded
	To           *time.Time
}

// EventResponse is one row of the list/export — the stored event plus the actor's
// display name joined in from "users". bson tags let the aggregation $project
// decode straight into this struct.
type EventResponse struct {
	ID         string         `bson:"_id" json:"id"`
	Action     string         `bson:"action" json:"action"`
	ActorID    string         `bson:"actor_id" json:"actor_id,omitempty"`
	ActorName  string         `bson:"actor_name" json:"actor_name,omitempty"`
	TargetType string         `bson:"target_type" json:"target_type,omitempty"`
	TargetID   string         `bson:"target_id" json:"target_id,omitempty"`
	TargetName string         `bson:"target_name" json:"target_name,omitempty"`
	Details    map[string]any `bson:"details" json:"details,omitempty"`
	DeviceID   string         `bson:"device_id" json:"device_id,omitempty"`
	IP         string         `bson:"ip" json:"ip,omitempty"`
	CreatedAt  time.Time      `bson:"created_at" json:"created_at"`
}
