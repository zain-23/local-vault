package api

import "net/http"

type Workspace struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type WorkspaceMembership struct {
	Workspace Workspace `json:"workspace"`
	Role      string    `json:"role"`
}

func (c *Client) ListWorkspaces() ([]WorkspaceMembership, error) {
	var out []WorkspaceMembership
	if err := c.do(http.MethodGet, "/api/v1/workspaces", nil, &out, true); err != nil {
		return nil, err
	}
	return out, nil
}
