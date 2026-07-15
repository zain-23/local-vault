package member

import "time"

// InviteRequest is the POST /members/invite body. `oneof` blocks inviting an owner.
type InviteRequest struct {
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required,oneof=admin member"`
}

// JoinRequest is the POST /members/join body — just the raw token from the email link.
type JoinRequest struct {
	Token string `json:"token" validate:"required"`
}

// ChangeRoleRequest is the PUT /members/:userId/role body. `oneof` blocks promoting to owner.
type ChangeRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=admin member"`
}

// MemberResponse is one row of the members list — membership joined with user display fields.
type MemberResponse struct {
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	Role      string    `json:"role"`
	JoinedAt  time.Time `json:"joined_at"`
}

// InviteResponse is one pending invite — never includes the token.
type InviteResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	InvitedBy string    `json:"invited_by"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// JoinResponse tells the client which workspace they just joined and with what role.
type JoinResponse struct {
	WorkspaceID string `json:"workspace_id"`
	Role        string `json:"role"`
}
