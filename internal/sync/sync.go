package sync

// sync.go handles vault encryption for peer transfer
// Uses X25519 Diffie-Hellman key exchange + AES-256-GCM encryption
//
// Why X25519?
// Both peers derive the SAME shared secret independently:
// A_private * B_public = B_private * A_public
// This is mathematically guaranteed — no server needed

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"golang.org/x/crypto/curve25519"
)

// VaultPayload is what travels between peers
// Contains encrypted secrets + metadata
type VaultPayload struct {
	Version string    `json:"version"`
	Secrets []byte    `json:"secrets"` // AES-256-GCM encrypted secrets JSON
	Nonce   []byte    `json:"nonce"`   // AES-GCM nonce
	SentAt  time.Time `json:"sent_at"`
	FromID  string    `json:"from_id"` // sender device ID
}

// SecretEntry is one secret in the transfer
type SecretEntry struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	Env       string    `json:"env"`
	UpdatedAt time.Time `json:"updated_at"`
}

// EncryptForPeer encrypts secrets so only the recipient can decrypt
//
// Process:
// 1. X25519 key exchange → shared secret
// 2. SHA256(shared secret) → AES key
// 3. AES-256-GCM encrypt secrets JSON
//
// senderX25519Private  → our 32 byte X25519 private key
// recipientX25519Public → their 32 byte X25519 public key
func EncryptForPeer(
	secrets []SecretEntry,
	senderX25519Private []byte,
	recipientX25519Public []byte,
	senderDeviceID string,
) ([]byte, error) {
	// Serialize secrets to JSON bytes
	secretsJSON, err := json.Marshal(secrets)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize secrets: %w", err)
	}

	// X25519 Diffie-Hellman key exchange
	// sharedSecret = our_private * their_public
	// Recipient computes: their_private * our_public = same result
	sharedSecret, err := curve25519.X25519(senderX25519Private, recipientX25519Public)
	if err != nil {
		return nil, fmt.Errorf("X25519 key exchange failed: %w", err)
	}

	// Hash shared secret to get 32 byte AES key
	// SHA256 provides domain separation and uniform distribution
	aesKey := sha256.Sum256(sharedSecret)

	// AES-256-GCM encryption
	block, err := aes.NewCipher(aesKey[:])
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Random nonce — must be unique per encryption
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt and authenticate
	encryptedSecrets := gcm.Seal(nil, nonce, secretsJSON, nil)

	// Build final payload
	payload := VaultPayload{
		Version: "1.0.0",
		Secrets: encryptedSecrets,
		Nonce:   nonce,
		SentAt:  time.Now(),
		FromID:  senderDeviceID,
	}

	return json.Marshal(payload)
}

// DecryptFromPeer decrypts secrets received from a peer
//
// recipientX25519Private → our 32 byte X25519 private key
// senderX25519Public     → sender's 32 byte X25519 public key
func DecryptFromPeer(
	payloadBytes []byte,
	recipientX25519Private []byte,
	senderX25519Public []byte,
) ([]SecretEntry, error) {
	// Parse outer payload envelope
	var payload VaultPayload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse payload: %w", err)
	}

	// X25519 key exchange — produces SAME shared secret as sender
	// Math: recipient_private * sender_public = sender_private * recipient_public
	// This equality is guaranteed by X25519 mathematics
	sharedSecret, err := curve25519.X25519(recipientX25519Private, senderX25519Public)
	if err != nil {
		return nil, fmt.Errorf("X25519 key exchange failed: %w", err)
	}

	// Derive same AES key
	aesKey := sha256.Sum256(sharedSecret)

	// AES-256-GCM decryption
	block, err := aes.NewCipher(aesKey[:])
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Decrypt and verify authentication tag
	// Returns error if data was tampered with
	secretsJSON, err := gcm.Open(nil, payload.Nonce, payload.Secrets, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed — keys do not match or data corrupted")
	}

	// Parse decrypted JSON into secrets
	var secrets []SecretEntry
	if err := json.Unmarshal(secretsJSON, &secrets); err != nil {
		return nil, fmt.Errorf("failed to parse decrypted secrets: %w", err)
	}

	return secrets, nil
}
