package api

import (
	"fmt"
	"net/http"
	"time"
)

type InviteCollaboratorRequest struct {
	Email      string `json:"email"`
	DeviceID   string `json:"device_id"`
	Code       string `json:"code"`
	WrappedDEK []byte `json:"wrapped_dek"`
}

type Collaborator struct {
	ID        string    `json:"id"`
	VaultID   string    `json:"vault_id"`
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	InvitedBy string    `json:"invited_by"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type JoinByCodeRequest struct {
	Code            string `json:"code"`
	DeviceID        string `json:"device_id"`
	DeviceName      string `json:"device_name"`
	PublicKey       []byte `json:"public_key"`
	X25519PublicKey []byte `json:"x25519_public_key"`
}

func (c *Client) InviteCollaborator(workspaceID, vaultID string, req InviteCollaboratorRequest) (*Collaborator, error) {
	var out Collaborator
	path := fmt.Sprintf("/api/v1/workspaces/%s/vaults/%s/collaborators", workspaceID, vaultID)
	if err := c.do(http.MethodPost, path, req, &out, true); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ListCollaborators(workspaceID, vaultID string) ([]Collaborator, error) {
	var out []Collaborator
	path := fmt.Sprintf("/api/v1/workspaces/%s/vaults/%s/collaborators", workspaceID, vaultID)
	if err := c.do(http.MethodGet, path, nil, &out, true); err != nil {
		return nil, err
	}
	if out == nil {
		out = []Collaborator{}
	}
	return out, nil
}

func (c *Client) RevokeCollaborator(workspaceID, vaultID, collabID string) error {
	path := fmt.Sprintf("/api/v1/workspaces/%s/vaults/%s/collaborators/%s", workspaceID, vaultID, collabID)
	return c.do(http.MethodDelete, path, nil, nil, true)
}

func (c *Client) JoinByCode(req JoinByCodeRequest) (*JoinResponse, error) {
	var out JoinResponse
	if err := c.do(http.MethodPost, "/api/v1/join-code", req, &out, true); err != nil {
		return nil, err
	}
	return &out, nil
}
