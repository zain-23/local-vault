// Mirrors server/internal/device/dto.go. Timestamps arrive as ISO strings.

export type DeviceDecisionAction = "approve" | "deny";

// GET /device/authorize/:userCode — what the browser shows before deciding.
// `status` is the request's lifecycle: "pending" until the user acts, then
// "approved" / "denied", or "expired" once the TTL passes.
export interface ApprovalDetails {
  device_name: string;
  ip: string;
  status: string;
  created_at: string;
  expires_at: string;
}

// PUT /device/authorize/:userCode body.
export interface DecisionInput {
  action: DeviceDecisionAction;
}

// GET /device — one authorized device. List/revoke UI is future work; the type
// lives here so the service and that screen agree when it lands.
export interface Device {
  id: string;
  name: string;
  ip: string;
  last_seen_at: string;
  authorized_at: string;
}
