package account

import "time"

// ProfileResponse — GET /account/me (secrets excluded).
type ProfileResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	Onboarded bool      `json:"onboarded"`
	CreatedAt time.Time `json:"created_at"`
}

// UpdateProfileRequest — PUT /account/me. All optional; a nil/empty field means
// "leave unchanged". Onboarded is a pointer so an explicit false is distinct
// from "not sent" — the onboarding flow sets it true, and a name edit must
// never reset it just because the flag was omitted. AvatarURL is intentionally
// not client-settable here — it's server-managed, populated only from GitHub on login.
type UpdateProfileRequest struct {
	Name      string `json:"name" validate:"omitempty,min=2,max=50"`
	Onboarded *bool  `json:"onboarded"`
}

// SessionResponse — one row of GET /account/sessions.
type SessionResponse struct {
	ID        string    `json:"id"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	Current   bool      `json:"current"` // true for the session making this request
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
