package workspace

type CreateWorkspaceRequest struct {
	Name string `json:"name" validate:"required,min=2,max=50"`
}

// UpdateWorkspaceRequest is the PUT /workspaces/:id body — rename only for now
type UpdateWorkspaceRequest struct {
	Name string `json:"name" validate:"required,min=2,max=50"`
}

type WorkspaceResponse struct {
	Workspace Workspace `json:"workspace"`
	Role      string    `json:"role"` // caller's role: owner | admin | member
}
