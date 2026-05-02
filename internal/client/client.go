package client

// client.go handles all HTTP communication with signaling server

import (
	"bytes"         // bytes.NewBuffer — wraps byte slice for HTTP body
	"encoding/json" // JSON encode/decode
	"fmt"
	"io"       // io.ReadAll — reads HTTP response body
	"net/http" // standard Go HTTP client
	"time"
)

// Client wraps HTTP communication with signaling server
type Client struct {
	baseURL    string       // signaling server URL
	httpClient *http.Client // underlying HTTP client
	deviceID   string       // this device's ID
}

// New creates a new signaling server client
func New(baseURL, deviceID string) *Client {
	return &Client{
		baseURL:  baseURL,
		deviceID: deviceID,
		httpClient: &http.Client{
			// Timeout prevents hanging forever
			// Like: axios.defaults.timeout = 10000 in JS
			Timeout: 10 * time.Second,
		},
	}
}

// ===== REQUEST/RESPONSE TYPES =====

// CreateInviteRequest is sent to POST /invite
type CreateInviteRequest struct {
	Code      string `json:"code"`
	DeviceID  string `json:"device_id"`
	PublicKey []byte `json:"public_key"`
	IPHint    string `json:"ip_hint"`
}

// CreateInviteResponse is returned from POST /invite
type CreateInviteResponse struct {
	Success   bool      `json:"success"`
	ExpiresAt time.Time `json:"expires_at"`
}

// PeerInfo is returned from GET /invite/:code
type PeerInfo struct {
	DeviceID  string `json:"device_id"`
	PublicKey []byte `json:"public_key"`
	IPHint    string `json:"ip_hint"`
}

// SendMessageRequest is sent to POST /messages
type SendMessageRequest struct {
	ForDeviceID   string `json:"for_device_id"`
	FromDeviceID  string `json:"from_device_id"`
	FromPublicKey []byte `json:"from_public_key"` // ← ADD THIS
	Payload       []byte `json:"payload"`
}

// PendingMessage is returned from GET /messages/:deviceID
type PendingMessage struct {
	ID            string    `json:"id"`
	ForDeviceID   string    `json:"for_device_id"`
	FromDeviceID  string    `json:"from_device_id"`
	FromPublicKey []byte    `json:"from_public_key"` // ← ADD THIS
	Payload       []byte    `json:"payload"`
	CreatedAt     time.Time `json:"created_at"`
}

// MessagesResponse is returned from GET /messages/:deviceID
type MessagesResponse struct {
	Messages []*PendingMessage `json:"messages"`
	Count    int               `json:"count"`
}

// ===== API METHODS =====

// HealthCheck verifies signaling server is reachable
// Called before any operation that needs the server
func (c *Client) HealthCheck() error {
	resp, err := c.httpClient.Get(c.baseURL + "/health")
	if err != nil {
		return fmt.Errorf("cannot reach signaling server at %s\n  Make sure server is running: go run server/main.go", c.baseURL)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("signaling server returned status %d", resp.StatusCode)
	}

	return nil
}

// CreateInvite registers an invite code on the server
// Called by: lv invite
func (c *Client) CreateInvite(req CreateInviteRequest) (*CreateInviteResponse, error) {
	// Encode request to JSON
	// Like: JSON.stringify(req) in JS
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	// Make POST request
	// bytes.NewBuffer wraps body bytes for http.Post
	resp, err := c.httpClient.Post(
		c.baseURL+"/invite",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create invite: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("server error: %s", string(respBody))
	}

	// Decode JSON response into struct
	var result CreateInviteResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// LookupInvite retrieves peer info using an invite code
// Called by: lv join LV-XXXX
// Returns Dev A's public key and IP hint
func (c *Client) LookupInvite(code string) (*PeerInfo, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/invite/" + code)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup invite: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("invite code not found or expired — ask your teammate to run: lv invite")
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("server error: %s", string(respBody))
	}

	var peer PeerInfo
	if err := json.Unmarshal(respBody, &peer); err != nil {
		return nil, err
	}

	return &peer, nil
}

// SendMessage stores encrypted message for offline peer
// Called when syncing changes to an offline peer
func (c *Client) SendMessage(req SendMessageRequest) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Post(
		c.baseURL+"/messages",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server error: %s", string(respBody))
	}

	return nil
}

// GetMessages retrieves pending messages for this device
// Called by: lv sync
// Returns all encrypted messages stored while device was offline
func (c *Client) GetMessages() (*MessagesResponse, error) {
	resp, err := c.httpClient.Get(
		fmt.Sprintf("%s/messages/%s", c.baseURL, c.deviceID),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("server error: %s", string(respBody))
	}

	var result MessagesResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
