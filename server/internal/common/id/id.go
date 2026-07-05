package id

import (
	"crypto/rand"
	"encoding/hex"
)

// Generate creates a random ID with prefix like "usr_a1b2c3..." — prefix tells you what type it is at a glance
func Generate(prefix string, length int) string {
	// make creates a byte slice (like an array) of given size
	b := make([]byte, length)
	// crypto/rand fills with secure random bytes — not math/rand which is predictable
	rand.Read(b)
	// hex converts bytes to readable characters: [0xA1] → "a1"
	return prefix + hex.EncodeToString(b)
}
