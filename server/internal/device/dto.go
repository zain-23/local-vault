package device

import "time"

// AuthorizeRequest is the POST /device/authorize body — sent by the CLI, unauthenticated.
type AuthorizeRequest struct {
	DeviceName        string `json:"device_name" validate:"required,max=100"`
	DeviceFingerprint string `json:"device_fingerprint" validate:"required,max=200"`
}

// AuthorizeResponse is what the CLI gets back. device_code is secret (poll with it);
// user_code is the short code shown in the terminal and used in the browser URL.
type AuthorizeResponse struct {
	DeviceCode string `json:"device_code"`
	UserCode   string `json:"user_code"`
	VerifyURL  string `json:"verify_url"`
	ExpiresIn  int    `json:"expires_in"` // seconds until the request dies
	Interval   int    `json:"interval"`   // seconds the CLI should wait between polls
}

// ApprovalDetailsResponse is what the browser shows on the approval screen:
type ApprovalDetailsResponse struct {
	DeviceName string    `json:"device_name"`
	IP         string    `json:"ip"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiresAt  time.Time `json:"expires_at"`
}

// DecisionRequest is the PUT /device/authorize/:userCode body.
// WorkspaceID is required to approve and ignored on deny — hence omitempty, not required.
type DecisionRequest struct {
	Action      string `json:"action" validate:"required,oneof=approve deny"`
	WorkspaceID string `json:"workspace_id" validate:"omitempty"`
}

// PollRequest is the POST /device/authorize/poll body. The secret goes in the body,
// never the URL — URLs land in access logs, proxy history and referrer headers.
type PollRequest struct {
	DeviceCode string `json:"device_code" validate:"required"`
}

// PollResponse carries the status, plus tokens only once approved.
// omitempty means a pending poll is just {"status":"pending"}.
type PollResponse struct {
	Status       string `json:"status"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	WorkspaceID  string `json:"workspace_id,omitempty"`
	DeviceID     string `json:"device_id,omitempty"`
}

// DeviceResponse is one row of GET /device.
type DeviceResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	IP           string    `json:"ip"`
	LastSeenAt   time.Time `json:"last_seen_at"`
	AuthorizedAt time.Time `json:"authorized_at"`
}

// ListDevicesQuery is the GET /device query string — devices are always workspace-scoped.
type ListDevicesQuery struct {
	WorkspaceID string `query:"workspace_id" validate:"required"`
}
