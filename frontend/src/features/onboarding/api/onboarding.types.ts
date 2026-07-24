// A workspace as returned by the server (server/internal/workspace/model.go).
export interface Workspace {
  id: string;
  name: string;
  slug: string;
  owner_id: string;
  created_at: string;
  updated_at: string;
}

// POST /workspaces/ envelope data — the created workspace plus the caller's role.
export interface WorkspaceResponse {
  workspace: Workspace;
  role: string; // owner | admin | member — creator is always "owner"
}

// An authorized device/terminal (server device.DeviceResponse). Its presence in
// the list is how onboarding knows a terminal finished `lv login`.
export interface Device {
  id: string;
  name: string;
  ip: string;
  last_seen_at: string;
  authorized_at: string;
}

// ---- request payloads ----

// Step 1 — the workspace name the user types in.
export interface CreateWorkspaceInput {
  name: string;
}

// Step 1 (rename path) — PUT /workspaces/:id body. Same shape as create, but
// applied to a workspace that already exists (e.g. after a page refresh).
export interface UpdateWorkspaceInput {
  name: string;
}

// Final step — mark the account as onboarded via PUT /account/me. The backend
// UpdateProfile only applies non-empty fields, so this is a targeted flip.
export interface CompleteOnboardingInput {
  onboarded: true;
}
