package workspace

type CreateWorkspaceRequest struct {
	Name string `json:"name" validate:"required,min=2,max=50"`
	Icon string `json:"icon" validate:"omitempty,oneof=vault lock key shield folder rocket wrench database cloud terminal boxes fingerprint"`
}

// UpdateWorkspaceRequest is the PUT /workspaces/:id body — name and/or icon
type UpdateWorkspaceRequest struct {
	Name string `json:"name" validate:"required,min=2,max=50"`
	Icon string `json:"icon" validate:"omitempty,oneof=vault lock key shield folder rocket wrench database cloud terminal boxes fingerprint"`
}

type WorkspaceResponse struct {
	Workspace Workspace `json:"workspace"`
	Role      string    `json:"role"` // caller's role: owner | admin | member
}
