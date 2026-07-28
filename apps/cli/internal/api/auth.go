package api

import "net/http"

// DeviceAuthorizeResponse is the server's reply to POST /device/authorize.
type DeviceAuthorizeResponse struct {
	DeviceCode string `json:"device_code"`
	UserCode   string `json:"user_code"`
	VerifyURL  string `json:"verify_url"`
	ExpiresIn  int    `json:"expires_in"`
	Interval   int    `json:"interval"`
}

// DeviceAuthorize starts the device login flow (unauthenticated).
func (c *Client) DeviceAuthorize(name, fingerprint string) (*DeviceAuthorizeResponse, error) {
	body := map[string]string{"device_name": name, "device_fingerprint": fingerprint}
	var res DeviceAuthorizeResponse
	if err := c.do(http.MethodPost, "/api/v1/device/authorize", body, &res, false); err != nil {
		return nil, err
	}
	return &res, nil
}

// DevicePollResponse carries the poll status and, once approved, the tokens.
type DevicePollResponse struct {
	Status       string `json:"status"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	DeviceID     string `json:"device_id"`
}

// DevicePoll checks the authorization status (unauthenticated; device_code is the credential).
func (c *Client) DevicePoll(deviceCode string) (*DevicePollResponse, error) {
	body := map[string]string{"device_code": deviceCode}
	var res DevicePollResponse
	if err := c.do(http.MethodPost, "/api/v1/device/authorize/poll", body, &res, false); err != nil {
		return nil, err
	}
	return &res, nil
}
