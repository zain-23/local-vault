package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	deviceID   string
	httpClient *http.Client
}

func New(baseURL, deviceID string) *Client {
	return &Client{
		baseURL:  baseURL,
		deviceID: deviceID,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// ===== REQUEST / RESPONSE TYPES =====

// Peer mirrors server Peer struct
type Peer struct {
	DeviceID        string    `json:"device_id"`
	DeviceName      string    `json:"device_name"`
	PublicKey       []byte    `json:"public_key"`
	X25519PublicKey []byte    `json:"x25519_public_key"`
	JoinedAt        time.Time `json:"joined_at"`
}

// RegisterVaultRequest sent on lv init
type RegisterVaultRequest struct {
	OwnerID         string `json:"owner_id"`
	OwnerName       string `json:"owner_name"`
	PublicKey       []byte `json:"public_key"`
	X25519PublicKey []byte `json:"x25519_public_key"`
}

// RegisterVaultResponse returned after vault created
type RegisterVaultResponse struct {
	VaultID   string    `json:"vault_id"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateTokenRequest sent on lv invite
type CreateTokenRequest struct {
	DeviceID  string     `json:"device_id"`
	Name      string     `json:"name"`
	ExpiresAt *time.Time `json:"expires_at"`
}

// Token returned from server
type Token struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at"`
}

// JoinRequest sent on lv join
type JoinRequest struct {
	Token           string `json:"token"`
	DeviceID        string `json:"device_id"`
	DeviceName      string `json:"device_name"`
	PublicKey       []byte `json:"public_key"`
	X25519PublicKey []byte `json:"x25519_public_key"`
}

// JoinResponse returned after joining
type JoinResponse struct {
	VaultID  string `json:"vault_id"`
	Snapshot []byte `json:"snapshot"` // nil if no snapshot yet
	Peers    []Peer `json:"peers"`
}

// PendingMessage for offline delivery
type PendingMessage struct {
	ID            string    `json:"id"`
	ForDeviceID   string    `json:"for_device_id"`
	FromDeviceID  string    `json:"from_device_id"`
	FromPublicKey []byte    `json:"from_public_key"`
	Payload       []byte    `json:"payload"`
	CreatedAt     time.Time `json:"created_at"`
}

// SendMessageRequest for offline delivery
type SendMessageRequest struct {
	ForDeviceID   string `json:"for_device_id"`
	FromDeviceID  string `json:"from_device_id"`
	FromPublicKey []byte `json:"from_public_key"`
	Payload       []byte `json:"payload"`
}

// MessagesResponse returned from GET /messages
type MessagesResponse struct {
	Messages []*PendingMessage `json:"messages"`
	Count    int               `json:"count"`
}

// ===== API METHODS =====

// HealthCheck verifies server is reachable
func (c *Client) HealthCheck() error {
	resp, err := c.httpClient.Get(c.baseURL + "/health")
	if err != nil {
		return fmt.Errorf(
			"cannot reach server at %s\n  Is it running?",
			c.baseURL,
		)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("server returned status %d", resp.StatusCode)
	}
	return nil
}

// RegisterVault creates vault on server during lv init
func (c *Client) RegisterVault(req RegisterVaultRequest) (*RegisterVaultResponse, error) {
	var result RegisterVaultResponse
	err := c.post("/vaults", req, &result)
	return &result, err
}

// UploadSnapshot uploads encrypted vault snapshot
// Called by lv push
func (c *Client) UploadSnapshot(vaultID string, snapshot []byte) error {
	body := map[string]interface{}{
		"device_id": c.deviceID,
		"snapshot":  snapshot,
	}
	return c.put(fmt.Sprintf("/vaults/%s/snapshot", vaultID), body, nil)
}

// DownloadSnapshot gets latest vault snapshot
// Called by lv sync
func (c *Client) DownloadSnapshot(vaultID string) ([]byte, time.Time, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/vaults/%s/snapshot", c.baseURL, vaultID),
		nil,
	)
	if err != nil {
		return nil, time.Time{}, err
	}

	// Send device ID so server can verify we are still a peer
	req.Header.Set("X-Device-ID", c.deviceID)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, time.Time{}, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == 403 {
		return nil, time.Time{}, fmt.Errorf(
			"access revoked — you are no longer in this vault",
		)
	}
	if resp.StatusCode == 404 {
		return nil, time.Time{}, fmt.Errorf(
			"no snapshot yet — run lv push first",
		)
	}
	if resp.StatusCode != 200 {
		return nil, time.Time{}, fmt.Errorf(
			"server error: %s", string(body),
		)
	}

	var result struct {
		Snapshot  []byte    `json:"snapshot"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	json.Unmarshal(body, &result)
	return result.Snapshot, result.UpdatedAt, nil
}

// GetVaultPeers returns all peers in the vault
func (c *Client) GetVaultPeers(vaultID string) ([]Peer, error) {
	resp, err := c.httpClient.Get(
		fmt.Sprintf("%s/vaults/%s", c.baseURL, vaultID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Peers []Peer `json:"peers"`
	}
	json.Unmarshal(body, &result)
	return result.Peers, nil
}

// CreateToken generates a join token for a teammate
// Called by lv invite
func (c *Client) CreateToken(vaultID string, req CreateTokenRequest) (*Token, error) {
	var result Token
	err := c.post(fmt.Sprintf("/vaults/%s/tokens", vaultID), req, &result)
	return &result, err
}

// ListTokens returns all active join tokens
// Called by lv invite --list
func (c *Client) ListTokens(vaultID string) ([]Token, error) {
	resp, err := c.httpClient.Get(
		fmt.Sprintf("%s/vaults/%s/tokens", c.baseURL, vaultID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Tokens []Token `json:"tokens"`
	}
	json.Unmarshal(body, &result)
	return result.Tokens, nil
}

// RevokeToken revokes a join token
// Called by lv invite --revoke TOKEN
func (c *Client) RevokeToken(vaultID, tokenID string) error {
	req, _ := http.NewRequest(
		"DELETE",
		fmt.Sprintf("%s/vaults/%s/tokens/%s", c.baseURL, vaultID, tokenID),
		nil,
	)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to revoke token")
	}
	return nil
}

// JoinVault joins vault using token
// Called by lv join TOKEN
func (c *Client) JoinVault(req JoinRequest) (*JoinResponse, error) {
	var result JoinResponse
	err := c.post("/join", req, &result)
	return &result, err
}

// SendMessage stores message for offline peer
func (c *Client) SendMessage(req SendMessageRequest) error {
	return c.post("/messages", req, nil)
}

// GetMessages retrieves pending messages
func (c *Client) GetMessages() (*MessagesResponse, error) {
	resp, err := c.httpClient.Get(
		fmt.Sprintf("%s/messages/%s", c.baseURL, c.deviceID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result MessagesResponse
	json.Unmarshal(body, &result)
	return &result, nil
}

// ===== HELPERS =====

func (c *Client) post(path string, body interface{}, result interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Post(
		c.baseURL+path,
		"application/json",
		bytes.NewBuffer(data),
	)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		var errResp struct {
			Error string `json:"error"`
		}
		json.Unmarshal(respBody, &errResp)
		return fmt.Errorf("%s", errResp.Error)
	}

	if result != nil {
		json.Unmarshal(respBody, result)
	}
	return nil
}

func (c *Client) put(path string, body interface{}, result interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, _ := http.NewRequest(
		"PUT",
		c.baseURL+path,
		bytes.NewBuffer(data),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		var errResp struct {
			Error string `json:"error"`
		}
		json.Unmarshal(respBody, &errResp)
		return fmt.Errorf("%s", errResp.Error)
	}

	if result != nil {
		json.Unmarshal(respBody, result)
	}
	return nil
}

// RemovePeer removes a peer from vault on server
// Called by lv revoke
func (c *Client) RemovePeer(vaultID, deviceID string) error {
	req, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("%s/vaults/%s/peers/%s",
			c.baseURL, vaultID, deviceID),
		nil,
	)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to remove peer: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("server failed to remove peer")
	}

	return nil
}
