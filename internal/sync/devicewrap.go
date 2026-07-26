package sync

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"

	"golang.org/x/crypto/curve25519"
)

const deviceWrapPrefix = "lv-device-wrap-v1"

// WrapDEKForDevice seals the vault DEK for a recipient device using X25519 + AES-GCM.
// Returns ciphertext (nonce||ct) suitable for server storage.
func WrapDEKForDevice(senderX25519Private, recipientX25519Public, dek []byte) ([]byte, error) {
	shared, err := curve25519.X25519(senderX25519Private, recipientX25519Public)
	if err != nil {
		return nil, fmt.Errorf("X25519 key exchange failed: %w", err)
	}
	aesKey := sha256.Sum256(append([]byte(deviceWrapPrefix), shared...))
	block, err := aes.NewCipher(aesKey[:])
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}
	return gcm.Seal(nonce, nonce, dek, nil), nil
}

// UnwrapDEKForDevice reverses WrapDEKForDevice.
func UnwrapDEKForDevice(recipientX25519Private, senderX25519Public, blob []byte) ([]byte, error) {
	shared, err := curve25519.X25519(recipientX25519Private, senderX25519Public)
	if err != nil {
		return nil, fmt.Errorf("X25519 key exchange failed: %w", err)
	}
	aesKey := sha256.Sum256(append([]byte(deviceWrapPrefix), shared...))
	block, err := aes.NewCipher(aesKey[:])
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(blob) < gcm.NonceSize() {
		return nil, fmt.Errorf("wrapped key too short")
	}
	nonce, ct := blob[:gcm.NonceSize()], blob[gcm.NonceSize():]
	dek, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return nil, fmt.Errorf("could not unwrap vault key")
	}
	return dek, nil
}
