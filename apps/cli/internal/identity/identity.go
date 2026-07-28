package identity

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/curve25519"
)

// Identity represents this device's public identity
// Stored in .lv/identity.json
type Identity struct {
	DeviceID        string    `json:"device_id"`
	DeviceName      string    `json:"device_name"`
	CreatedAt       time.Time `json:"created_at"`
	PublicKey       []byte    `json:"public_key"`        // Ed25519 — for signing
	X25519PublicKey []byte    `json:"x25519_public_key"` // X25519 — for encryption
}

// FullIdentity holds everything including private keys
// Private keys never leave this struct — never stored in JSON
type FullIdentity struct {
	Identity
	PrivateKey       ed25519.PrivateKey // Ed25519 private key for signing
	X25519PrivateKey []byte             // X25519 private key for encryption
}

const (
	identityFile = "identity.json"
	privateFile  = "identity.key"
	publicFile   = "identity.pub"
)

// Generate creates a brand new identity for this device
// Called once during: lv init
func Generate(lvDir string) (*FullIdentity, error) {
	// Generate Ed25519 keypair for signing
	// Used to prove messages come from us
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Ed25519 keypair: %w", err)
	}

	// Generate X25519 keypair for encryption
	// Used for secure key exchange with peers
	// Separate from Ed25519 — correct cryptographic practice
	var x25519Private [32]byte
	if _, err := rand.Read(x25519Private[:]); err != nil {
		return nil, fmt.Errorf("failed to generate X25519 private key: %w", err)
	}

	// Clamp private key bits as per RFC 7748
	// Required for X25519 to work correctly
	x25519Private[0] &= 248
	x25519Private[31] &= 127
	x25519Private[31] |= 64

	// Derive X25519 public key from private key
	// curve25519.Basepoint is the standard generator point
	x25519Public, err := curve25519.X25519(x25519Private[:], curve25519.Basepoint)
	if err != nil {
		return nil, fmt.Errorf("failed to derive X25519 public key: %w", err)
	}

	// Get hostname for human readable device name
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown-device"
	}

	id := &FullIdentity{
		Identity: Identity{
			DeviceID:        uuid.New().String(),
			DeviceName:      hostname,
			CreatedAt:       time.Now(),
			PublicKey:       publicKey,
			X25519PublicKey: x25519Public,
		},
		PrivateKey:       privateKey,
		X25519PrivateKey: x25519Private[:],
	}

	if err := id.Save(lvDir); err != nil {
		return nil, err
	}

	return id, nil
}

// Save writes identity files to disk
// Saves 3 files:
// identity.json → public metadata (safe to commit)
// identity.key  → both private keys (never commit)
// identity.pub  → Ed25519 public key (safe to share)
func (id *FullIdentity) Save(lvDir string) error {
	// Save public identity as JSON
	identityData, err := json.MarshalIndent(id.Identity, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize identity: %w", err)
	}

	identityPath := filepath.Join(lvDir, identityFile)
	if err := os.WriteFile(identityPath, identityData, 0644); err != nil {
		return fmt.Errorf("failed to save identity: %w", err)
	}

	// Save private keys to identity.key
	// One file contains both Ed25519 and X25519 private keys
	// Protected with 0600 permissions — only owner can read
	privatePath := filepath.Join(lvDir, privateFile)
	f, err := os.OpenFile(privatePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to create private key file: %w", err)
	}
	defer f.Close()

	// Write Ed25519 private key as first PEM block
	ed25519Block := &pem.Block{
		Type:  "ED25519 PRIVATE KEY",
		Bytes: id.PrivateKey,
	}
	if err := pem.Encode(f, ed25519Block); err != nil {
		return fmt.Errorf("failed to write Ed25519 private key: %w", err)
	}

	// Write X25519 private key as second PEM block
	x25519Block := &pem.Block{
		Type:  "X25519 PRIVATE KEY",
		Bytes: id.X25519PrivateKey,
	}
	if err := pem.Encode(f, x25519Block); err != nil {
		return fmt.Errorf("failed to write X25519 private key: %w", err)
	}

	// Save Ed25519 public key to identity.pub
	publicBlock := &pem.Block{
		Type:  "ED25519 PUBLIC KEY",
		Bytes: id.PublicKey,
	}
	publicPath := filepath.Join(lvDir, publicFile)
	if err := os.WriteFile(publicPath, pem.EncodeToMemory(publicBlock), 0644); err != nil {
		return fmt.Errorf("failed to save public key: %w", err)
	}

	return nil
}

// Load reads identity from disk
// Called before any command that needs identity
func Load(lvDir string) (*FullIdentity, error) {
	// Read public identity JSON
	identityPath := filepath.Join(lvDir, identityFile)
	data, err := os.ReadFile(identityPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("no identity found — run 'lv init' first")
		}
		return nil, fmt.Errorf("failed to read identity: %w", err)
	}

	var id Identity
	if err := json.Unmarshal(data, &id); err != nil {
		return nil, fmt.Errorf("failed to parse identity: %w", err)
	}

	// Read private keys file
	privatePath := filepath.Join(lvDir, privateFile)
	privateData, err := os.ReadFile(privatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	// Decode first PEM block — Ed25519 private key
	ed25519Block, rest := pem.Decode(privateData)
	if ed25519Block == nil {
		return nil, errors.New("failed to decode Ed25519 private key")
	}

	// Decode second PEM block — X25519 private key
	x25519Block, _ := pem.Decode(rest)
	if x25519Block == nil {
		return nil, errors.New("failed to decode X25519 private key — try running lv init again")
	}

	return &FullIdentity{
		Identity:         id,
		PrivateKey:       ed25519.PrivateKey(ed25519Block.Bytes),
		X25519PrivateKey: x25519Block.Bytes,
	}, nil
}

// Sign creates a cryptographic signature for a message
// Used to prove a message came from this device
func (id *FullIdentity) Sign(message []byte) []byte {
	return ed25519.Sign(id.PrivateKey, message)
}

// Verify checks if a signature was made by a specific public key
func Verify(publicKey []byte, message []byte, signature []byte) bool {
	return ed25519.Verify(publicKey, message, signature)
}

// PublicKeyString returns Ed25519 public key as PEM string
func (id *FullIdentity) PublicKeyString() string {
	block := &pem.Block{
		Type:  "ED25519 PUBLIC KEY",
		Bytes: id.PublicKey,
	}
	return string(pem.EncodeToMemory(block))
}

// Exists checks if identity already exists in a directory
func Exists(lvDir string) bool {
	_, err := os.Stat(filepath.Join(lvDir, identityFile))
	return err == nil
}
