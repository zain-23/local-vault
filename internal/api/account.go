package api

import (
	"net/http"
	"time"
)

// Account mirrors the server's GET /account/me payload.
type Account struct {
	ID               string    `json:"id"`
	Email            string    `json:"email"`
	Name             string    `json:"name"`
	AvatarURL        string    `json:"avatar_url"`
	TwoFactorEnabled bool      `json:"two_factor_enabled"`
	CreatedAt        time.Time `json:"created_at"`
}

// Me returns the logged-in account (authenticated).
func (c *Client) Me() (*Account, error) {
	var res Account
	if err := c.do(http.MethodGet, "/api/v1/account/me", nil, &res, true); err != nil {
		return nil, err
	}
	return &res, nil
}
