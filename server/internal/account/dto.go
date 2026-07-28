package account

import "time"

// ProfileResponse — GET /account/me (secrets excluded).
type ProfileResponse struct {
	ID               string    `json:"id"`
	Email            string    `json:"email"`
	Name             string    `json:"name"`
	AvatarURL        string    `json:"avatar_url,omitempty"`
	TwoFactorEnabled bool      `json:"two_factor_enabled"`
	Onboarded        bool      `json:"onboarded"`
	CreatedAt        time.Time `json:"created_at"`
}

// UpdateProfileRequest — PUT /account/me. All optional; a nil/empty field means
// "leave unchanged". Onboarded is a pointer so an explicit false is distinct
// from "not sent" — the onboarding flow sets it true, and a name/avatar edit
// must never reset it just because the flag was omitted.
type UpdateProfileRequest struct {
	Name      string `json:"name" validate:"omitempty,min=2,max=50"`
	AvatarURL string `json:"avatar_url" validate:"omitempty,url"`
	Onboarded *bool  `json:"onboarded"`
}

// ChangePasswordRequest — PUT /account/password.
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=72"`
}

// Enable2FAResponse — the secret + otpauth URL to render as a QR (shown once).
type Enable2FAResponse struct {
	Secret     string `json:"secret"`
	OtpauthURL string `json:"otpauth_url"`
}

// Verify2FARequest / Response — confirm setup, receive backup codes (shown once).
type Verify2FARequest struct {
	TOTPCode string `json:"totp_code" validate:"required"`
}
type Verify2FAResponse struct {
	BackupCodes []string `json:"backup_codes"`
}

// Disable2FARequest — a TOTP code OR a backup code (at least one required, checked in the service).
type Disable2FARequest struct {
	TOTPCode   string `json:"totp_code" validate:"omitempty"`
	BackupCode string `json:"backup_code" validate:"omitempty"`
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
