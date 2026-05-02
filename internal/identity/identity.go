package identity

// identity.go manages device identity
// Every machine that uses LocalVault has:
// 1. A unique device ID (like a username)
// 2. A keypair (private + public key)
//    Private key = stays on this machine forever
//    Public key  = shared with teammates during invite

// We use Ed25519 for signing and X25519 for encryption
// Ed25519 → proves a message came from you (signing)
// X25519  → encrypts data so only recipient can read (key exchange)

import (
	"crypto/ed25519" // Ed25519 signing algorithm
	// same algorithm GitHub uses for SSH keys
	"crypto/rand"   // cryptographically secure random numbers
	"encoding/json" // JSON marshal/unmarshal
	"encoding/pem"  // PEM encoding (standard key file format)

	// same format as .pem files in SSL certificates
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	// generates unique IDs
	// like crypto.randomUUID() in JS
)

// Identity represents this device's identity
// Stored in .lv/identity.json
type Identity struct {
	DeviceID   string    `json:"device_id"`   // unique ID for this machine
	DeviceName string    `json:"device_name"` // human readable name e.g. "Arslan's MacBook"
	CreatedAt  time.Time `json:"created_at"`
	PublicKey  []byte    `json:"public_key"` // safe to share with teammates
}

// FullIdentity holds everything including private key
// Private key is NEVER stored in identity.json
// It lives in a separate file: identity.key
type FullIdentity struct {
	Identity
	PrivateKey ed25519.PrivateKey // kept only in memory after loading
}

// File names inside .lv/ folder
const (
	identityFile = "identity.json" // public info — safe to commit
	privateFile  = "identity.key"  // private key — NEVER commit
	publicFile   = "identity.pub"  // public key — safe to share
)

// Generate creates a brand new identity for this device
// Called once during: lv init
// Never called again on same machine
func Generate(lvDir string) (*FullIdentity, error) {
	// Generate Ed25519 keypair
	// ed25519.GenerateKey returns (publicKey, privateKey, error)
	// privateKey in Go's ed25519 actually contains both
	// private and public parts together (64 bytes total)
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate keypair: %w", err)
	}

	// Generate unique device ID
	// uuid.New() creates a v4 UUID like: "a3f9b2c1-4d5e-6f7a-8b9c-0d1e2f3a4b5c"
	deviceID := uuid.New().String()

	// Get computer hostname as device name
	// Like: os.hostname() in Node.js
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown-device" // fallback if hostname unavailable
	}

	// Build identity struct
	id := &FullIdentity{
		Identity: Identity{
			DeviceID:   deviceID,
			DeviceName: hostname,
			CreatedAt:  time.Now(),
			PublicKey:  publicKey,
		},
		PrivateKey: privateKey,
	}

	// Save to disk
	if err := id.Save(lvDir); err != nil {
		return nil, err
	}

	return id, nil
}

// Save writes identity files to disk
// Saves 3 files:
// identity.json → public metadata
// identity.key  → private key (protected)
// identity.pub  → public key (shareable)
func (id *FullIdentity) Save(lvDir string) error {
	// Save public identity JSON
	// This file is safe to read by anyone
	identityData, err := json.MarshalIndent(id.Identity, "", "  ")
	// MarshalIndent = JSON.stringify(obj, null, 2) in JS
	if err != nil {
		return fmt.Errorf("failed to serialize identity: %w", err)
	}

	identityPath := filepath.Join(lvDir, identityFile)
	if err := os.WriteFile(identityPath, identityData, 0644); err != nil {
		// 0644 = owner can read/write, others can read
		return fmt.Errorf("failed to save identity: %w", err)
	}

	// Save private key in PEM format
	// PEM is standard format for cryptographic keys
	// Same format as: -----BEGIN RSA PRIVATE KEY-----
	// We use 0600 permission = ONLY owner can read/write
	// No one else on the system can read this file
	privateBlock := &pem.Block{
		Type:  "ED25519 PRIVATE KEY",
		Bytes: id.PrivateKey, // raw private key bytes
	}

	privatePath := filepath.Join(lvDir, privateFile)
	privateFile, err := os.OpenFile(privatePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to create private key file: %w", err)
	}
	defer privateFile.Close()

	if err := pem.Encode(privateFile, privateBlock); err != nil {
		return fmt.Errorf("failed to write private key: %w", err)
	}

	// Save public key in PEM format
	// This file is safe to share with teammates
	// During invite, we send this to the other device
	publicBlock := &pem.Block{
		Type:  "ED25519 PUBLIC KEY",
		Bytes: id.PublicKey,
	}

	publicPath := filepath.Join(lvDir, publicFile)
	if err := os.WriteFile(
		publicPath,
		pem.EncodeToMemory(publicBlock), // convert PEM block to bytes
		0644,
	); err != nil {
		return fmt.Errorf("failed to save public key: %w", err)
	}

	return nil
}

// Load reads identity from disk
// Called at start of any command that needs identity (invite, join, sync)
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

	// Parse JSON into Identity struct
	var id Identity
	if err := json.Unmarshal(data, &id); err != nil {
		return nil, fmt.Errorf("failed to parse identity: %w", err)
	}

	// Read private key file
	privatePath := filepath.Join(lvDir, privateFile)
	privateData, err := os.ReadFile(privatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}

	// Decode PEM block
	// pem.Decode returns the first PEM block found in data
	block, _ := pem.Decode(privateData)
	if block == nil {
		return nil, errors.New("failed to decode private key — file may be corrupted")
	}

	// Convert raw bytes to ed25519.PrivateKey type
	// ed25519.PrivateKey is just []byte with a specific length (64 bytes)
	privateKey := ed25519.PrivateKey(block.Bytes)

	return &FullIdentity{
		Identity:   id,
		PrivateKey: privateKey,
	}, nil
}

// Sign creates a cryptographic signature for a message
// Used to prove a message came from this device
// Like a digital signature on a document
//
// Example use: when sending vault changes to peers
// We sign the changes so receiver knows it really came from us
func (id *FullIdentity) Sign(message []byte) []byte {
	// ed25519.Sign(privateKey, message) returns 64-byte signature
	return ed25519.Sign(id.PrivateKey, message)
}

// Verify checks if a signature was made by a specific public key
// Used to verify messages from peers are authentic
// Returns true if signature is valid, false if tampered
func Verify(publicKey []byte, message []byte, signature []byte) bool {
	// ed25519.Verify returns bool
	return ed25519.Verify(publicKey, message, signature)
}

// PublicKeyString returns public key as readable string
// Used when displaying device info to user
func (id *FullIdentity) PublicKeyString() string {
	block := &pem.Block{
		Type:  "ED25519 PUBLIC KEY",
		Bytes: id.PublicKey,
	}
	return string(pem.EncodeToMemory(block))
}

// Exists checks if identity already exists in a directory
// Used by lv init to avoid overwriting existing identity
func Exists(lvDir string) bool {
	identityPath := filepath.Join(lvDir, identityFile)
	_, err := os.Stat(identityPath)
	return err == nil // true if file exists
}
