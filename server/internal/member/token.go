package member

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"strings"
)

// generateInviteToken creates a URL-safe random token for the email link.
func generateInviteToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// URL-safe alphabet (-, _) so the token survives email links / query strings.
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// sha256Hex hashes a token — we store the hash, never the raw token (leaked DB can't reuse it).
func sha256Hex(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:]) // h[:] converts the fixed-size array to a slice
}

// normalizeEmail lowercases and trims an address so invite/login emails compare reliably.
func normalizeEmail(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
