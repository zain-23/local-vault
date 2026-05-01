package crypto

// This package handles ALL encryption and decryption
// No other package touches raw secret values directly
// Think of this as your security boundary

import (
	"crypto/aes"    // AES encryption algorithm (standard library)
	"crypto/cipher" // cipher modes like GCM (standard library)
	"crypto/rand"   // cryptographically secure random numbers
	"crypto/sha256" // SHA256 hashing
	"encoding/hex"  // convert bytes to hex strings
	"errors"        // create error values
	"io"            // io.ReadFull for reading random bytes

	"golang.org/x/crypto/argon2" // Argon2id — best password hashing algorithm
	// same one used by 1Password
)

// Salt is random bytes mixed into key derivation
// Makes it impossible to use precomputed rainbow tables
// We store this alongside the encrypted data
const saltSize = 16 // 16 bytes = 128 bits

// DeriveKey turns a passphrase (human password) into a
// 32-byte encryption key using Argon2id algorithm
//
// Why Argon2id?
// - Intentionally slow and memory-hard
// - Makes brute force attacks take years not seconds
// - Winner of the Password Hashing Competition (2015)
// - Used by 1Password, Bitwarden, etc.
//
// JS equivalent: await argon2.hash(passphrase)
func DeriveKey(passphrase string, salt []byte) []byte {
	return argon2.IDKey(
		[]byte(passphrase), // convert string to bytes
		salt,
		1,       // time cost: 1 iteration
		64*1024, // memory cost: 64MB RAM required
		4,       // parallelism: use 4 CPU threads
		32,      // output key length: 32 bytes = 256 bits
	)
}

// GenerateSalt creates random bytes for use as salt
// Called once when vault is first created
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, saltSize) // allocate 16 bytes
	_, err := io.ReadFull(rand.Reader, salt)
	// io.ReadFull ensures ALL bytes are filled with random data
	return salt, err
}

// Encrypt takes plain data + key, returns encrypted bytes
//
// We use AES-256-GCM because:
//   - AES-256: military grade, 256-bit key
//   - GCM mode: authenticated encryption
//     (detects if encrypted data was tampered with)
//
// JS equivalent: crypto.createCipheriv('aes-256-gcm', key, iv)
func Encrypt(data []byte, key []byte) ([]byte, error) {
	// Create AES cipher block using our 32-byte key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Wrap block cipher in GCM mode
	// GCM adds authentication on top of encryption
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Nonce = "number used once"
	// Must be unique for every encryption operation
	// Like an IV (initialization vector) in JS crypto
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Seal encrypts and authenticates the data
	// We prepend nonce to ciphertext so we can extract it during decryption
	// Final format: [nonce][ciphertext][auth_tag]
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// Decrypt reverses Encrypt — takes encrypted bytes + key, returns plain data
// Returns error if data was tampered with (GCM authentication)
func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	// Recreate the same AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Extract nonce from the beginning of ciphertext
	// Remember we prepended it during Encrypt
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Open decrypts AND verifies authentication tag
	// If anyone modified the ciphertext, this returns an error
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// HashPassphrase creates a fingerprint of the passphrase
// Used to verify passphrase is correct before attempting decryption
// Stored in vault config so we can tell user "wrong password" early
func HashPassphrase(passphrase string) string {
	hash := sha256.Sum256([]byte(passphrase))
	return hex.EncodeToString(hash[:])
}
