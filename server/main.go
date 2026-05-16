package main

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// ===== DATA STRUCTURES =====

// Vault represents a registered vault on the server
// Created when Dev A runs lv init
type Vault struct {
	ID        string    `json:"id"`       // unique vault ID
	OwnerID   string    `json:"owner_id"` // Dev A device ID
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Snapshot  []byte    `json:"snapshot"` // latest encrypted vault blob
	Peers     []Peer    `json:"peers"`    // all registered peers
	Tokens    []Token   `json:"tokens"`   // active join tokens
}

// Peer represents a device that has joined the vault
type Peer struct {
	DeviceID        string    `json:"device_id"`
	DeviceName      string    `json:"device_name"`
	PublicKey       []byte    `json:"public_key"`        // Ed25519
	X25519PublicKey []byte    `json:"x25519_public_key"` // for encryption
	JoinedAt        time.Time `json:"joined_at"`
}

// Token represents a join token
// Like a GitHub personal access token
type Token struct {
	ID         string     `json:"id"`   // public half: lv_join_<id>
	Name       string     `json:"name"` // human label e.g. "Ahmed"
	CreatedAt  time.Time  `json:"created_at"`
	ExpiresAt  *time.Time `json:"expires_at"`  // nil = never expires
	Revoked    bool       `json:"revoked"`
	WrappedDEK []byte     `json:"wrapped_dek"` // opaque: DEK wrapped with the token secret
	Verifier   string     `json:"verifier"`    // one-way hash of the token secret
}

// PendingMessage for offline peer delivery
type PendingMessage struct {
	ID            string    `json:"id"`
	ForDeviceID   string    `json:"for_device_id"`
	FromDeviceID  string    `json:"from_device_id"`
	FromPublicKey []byte    `json:"from_public_key"`
	Payload       []byte    `json:"payload"`
	CreatedAt     time.Time `json:"created_at"`
	ExpiresAt     time.Time `json:"expires_at"`
}

// ===== STORE =====

type Store struct {
	mu       sync.RWMutex
	vaults   map[string]*Vault            // vaultID → vault
	messages map[string][]*PendingMessage // deviceID → messages
}

func NewStore() *Store {
	return &Store{
		vaults:   make(map[string]*Vault),
		messages: make(map[string][]*PendingMessage),
	}
}

// generateID creates a random hex ID
func generateID(prefix string, length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return prefix + hex.EncodeToString(b)
}

// GetVault finds vault by ID
func (s *Store) GetVault(id string) (*Vault, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.vaults[id]
	return v, ok
}

// GetVaultByToken finds vault using a join token
func (s *Store) GetVaultByToken(token string) (*Vault, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, vault := range s.vaults {
		for _, t := range vault.Tokens {
			if t.ID == token && !t.Revoked {
				// Check expiry
				if t.ExpiresAt != nil && time.Now().After(*t.ExpiresAt) {
					continue
				}
				return vault, true
			}
		}
	}
	return nil, false
}

// Cleanup removes expired messages
func (s *Store) Cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for deviceID, msgs := range s.messages {
		var valid []*PendingMessage
		for _, msg := range msgs {
			if now.Before(msg.ExpiresAt) {
				valid = append(valid, msg)
			}
		}
		if len(valid) == 0 {
			delete(s.messages, deviceID)
		} else {
			s.messages[deviceID] = valid
		}
	}
}

// ===== MAIN =====

func main() {
	store := NewStore()

	// Cleanup goroutine
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for range ticker.C {
			store.Cleanup()
			log.Println("🧹 Cleaned up expired messages")
		}
	}()

	app := buildApp(store)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 LocalVault Server v2.0 running on port %s", port)
	app.Listen(":" + port)
}

// buildApp registers all routes on a fresh fiber app. Extracted from
// main() so the HTTP handlers can be exercised in tests via app.Test().
func buildApp(store *Store) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		},
	})

	app.Use(cors.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${method} ${path} → ${status}\n",
	}))

	// ===== ROUTES =====

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"version": "2.0.0",
			"time":    time.Now(),
		})
	})

	// ── Vault Routes ────────────────────────────────────────

	// POST /vaults — Register new vault (called by lv init)
	app.Post("/vaults", func(c *fiber.Ctx) error {
		var body struct {
			OwnerID         string `json:"owner_id"`
			OwnerName       string `json:"owner_name"`
			PublicKey       []byte `json:"public_key"`
			X25519PublicKey []byte `json:"x25519_public_key"`
		}

		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}

		store.mu.Lock()
		defer store.mu.Unlock()

		// Generate vault ID
		vaultID := generateID("vault_", 8)

		// Create vault with owner as first peer
		vault := &Vault{
			ID:        vaultID,
			OwnerID:   body.OwnerID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Peers: []Peer{
				{
					DeviceID:        body.OwnerID,
					DeviceName:      body.OwnerName,
					PublicKey:       body.PublicKey,
					X25519PublicKey: body.X25519PublicKey,
					JoinedAt:        time.Now(),
				},
			},
		}

		store.vaults[vaultID] = vault

		log.Printf("🔐 Vault created: %s by %s", vaultID, body.OwnerID)

		return c.JSON(fiber.Map{
			"vault_id":   vaultID,
			"created_at": vault.CreatedAt,
		})
	})

	// GET /vaults/:id — Get vault info and peer list
	app.Get("/vaults/:id", func(c *fiber.Ctx) error {
		vault, ok := store.GetVault(c.Params("id"))
		if !ok {
			return c.Status(404).JSON(fiber.Map{"error": "vault not found"})
		}

		return c.JSON(fiber.Map{
			"id":         vault.ID,
			"peers":      vault.Peers,
			"updated_at": vault.UpdatedAt,
		})
	})

	// PUT /vaults/:id/snapshot — Upload encrypted vault snapshot
	// Called by lv push
	app.Put("/vaults/:id/snapshot", func(c *fiber.Ctx) error {
		var body struct {
			DeviceID string `json:"device_id"`
			Snapshot []byte `json:"snapshot"` // encrypted blob
		}

		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}

		store.mu.Lock()
		defer store.mu.Unlock()

		vault, ok := store.vaults[c.Params("id")]
		if !ok {
			return c.Status(404).JSON(fiber.Map{"error": "vault not found"})
		}

		// Verify device is a peer
		isPeer := false
		for _, p := range vault.Peers {
			if p.DeviceID == body.DeviceID {
				isPeer = true
				break
			}
		}

		if !isPeer {
			return c.Status(403).JSON(fiber.Map{
				"error": "device is not a vault peer",
			})
		}

		// Update snapshot
		vault.Snapshot = body.Snapshot
		vault.UpdatedAt = time.Now()

		log.Printf("📤 Snapshot updated: %s by %s",
			c.Params("id"), body.DeviceID)

		return c.JSON(fiber.Map{
			"success":    true,
			"updated_at": vault.UpdatedAt,
		})
	})

	// GET /vaults/:id/snapshot — Download latest snapshot
	app.Get("/vaults/:id/snapshot", func(c *fiber.Ctx) error {
		// Get requesting device ID from header
		requestingDevice := c.Get("X-Device-ID")

		vault, ok := store.GetVault(c.Params("id"))
		if !ok {
			return c.Status(404).JSON(fiber.Map{
				"error": "vault not found",
			})
		}

		// Verify requester is still a peer
		if requestingDevice != "" {
			isPeer := false
			for _, p := range vault.Peers {
				if p.DeviceID == requestingDevice {
					isPeer = true
					break
				}
			}
			if !isPeer {
				return c.Status(403).JSON(fiber.Map{
					"error": "access revoked",
				})
			}
		}

		if vault.Snapshot == nil {
			return c.Status(404).JSON(fiber.Map{
				"error": "no snapshot yet",
			})
		}

		return c.JSON(fiber.Map{
			"snapshot":   vault.Snapshot,
			"updated_at": vault.UpdatedAt,
		})
	})
	// ── Token Routes ─────────────────────────────────────────

	// POST /vaults/:id/tokens — Generate join token
	// Called by lv invite
	app.Post("/vaults/:id/tokens", func(c *fiber.Ctx) error {
		var body struct {
			DeviceID   string     `json:"device_id"`  // must be owner
			Name       string     `json:"name"`       // e.g. "Ahmed"
			ExpiresAt  *time.Time `json:"expires_at"` // nil = never
			WrappedDEK []byte     `json:"wrapped_dek"`
			Verifier   string     `json:"verifier"`
		}

		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}

		store.mu.Lock()
		defer store.mu.Unlock()

		vault, ok := store.vaults[c.Params("id")]
		if !ok {
			return c.Status(404).JSON(fiber.Map{"error": "vault not found"})
		}

		// Generate token (public half only — the secret never reaches us)
		tokenID := generateID("lv_join_", 16)

		token := Token{
			ID:         tokenID,
			Name:       body.Name,
			CreatedAt:  time.Now(),
			ExpiresAt:  body.ExpiresAt,
			Revoked:    false,
			WrappedDEK: body.WrappedDEK,
			Verifier:   body.Verifier,
		}

		vault.Tokens = append(vault.Tokens, token)

		log.Printf("🎟️  Token created for %s in vault %s",
			body.Name, c.Params("id"))

		return c.JSON(fiber.Map{
			"id":         tokenID,
			"name":       body.Name,
			"created_at": time.Now(),
			"expires_at": body.ExpiresAt,
		})
	})

	// GET /vaults/:id/tokens — List all tokens
	// Called by lv invite --list
	app.Get("/vaults/:id/tokens", func(c *fiber.Ctx) error {
		vault, ok := store.GetVault(c.Params("id"))
		if !ok {
			return c.Status(404).JSON(fiber.Map{"error": "vault not found"})
		}

		// Return only active tokens
		var active []Token
		for _, t := range vault.Tokens {
			if !t.Revoked {
				active = append(active, t)
			}
		}

		return c.JSON(fiber.Map{"tokens": active})
	})

	// DELETE /vaults/:id/tokens/:tokenID — Revoke token
	// Called by lv invite --revoke TOKEN
	app.Delete("/vaults/:id/tokens/:tokenID", func(c *fiber.Ctx) error {
		store.mu.Lock()
		defer store.mu.Unlock()

		vault, ok := store.vaults[c.Params("id")]
		if !ok {
			return c.Status(404).JSON(fiber.Map{"error": "vault not found"})
		}

		tokenID := c.Params("tokenID")
		found := false
		for i, t := range vault.Tokens {
			if t.ID == tokenID {
				vault.Tokens[i].Revoked = true
				found = true
				break
			}
		}

		if !found {
			return c.Status(404).JSON(fiber.Map{"error": "token not found"})
		}

		log.Printf("🚫 Token revoked: %s", tokenID)

		return c.JSON(fiber.Map{"success": true})
	})

	// ── Join Route ───────────────────────────────────────────

	// POST /join — Join vault using token
	// Called by lv join TOKEN
	app.Post("/join", func(c *fiber.Ctx) error {
		var body struct {
			Token           string `json:"token"`
			Verifier        string `json:"verifier"`
			DeviceID        string `json:"device_id"`
			DeviceName      string `json:"device_name"`
			PublicKey       []byte `json:"public_key"`
			X25519PublicKey []byte `json:"x25519_public_key"`
		}

		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}

		store.mu.Lock()
		defer store.mu.Unlock()

		// Find vault + the specific token by its public id.
		var targetVault *Vault
		var matched *Token
		for _, vault := range store.vaults {
			for i := range vault.Tokens {
				t := &vault.Tokens[i]
				if t.ID == body.Token && !t.Revoked {
					if t.ExpiresAt == nil || time.Now().Before(*t.ExpiresAt) {
						targetVault = vault
						matched = t
						break
					}
				}
			}
			if targetVault != nil {
				break
			}
		}

		// A mismatched verifier is reported identically to an unknown
		// token — no oracle for whether the id exists.
		if targetVault == nil || matched == nil ||
			(matched.Verifier != "" && matched.Verifier != body.Verifier) {
			return c.Status(404).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		// Check if already a peer
		for _, p := range targetVault.Peers {
			if p.DeviceID == body.DeviceID {
				// Already joined — return vault info
				return c.JSON(fiber.Map{
					"vault_id":    targetVault.ID,
					"snapshot":    targetVault.Snapshot,
					"peers":       targetVault.Peers,
					"wrapped_dek": matched.WrappedDEK,
					"message":     "already a peer",
				})
			}
		}

		// Add as new peer
		newPeer := Peer{
			DeviceID:        body.DeviceID,
			DeviceName:      body.DeviceName,
			PublicKey:       body.PublicKey,
			X25519PublicKey: body.X25519PublicKey,
			JoinedAt:        time.Now(),
		}

		targetVault.Peers = append(targetVault.Peers, newPeer)
		targetVault.UpdatedAt = time.Now()

		log.Printf("👋 %s joined vault %s",
			body.DeviceName, targetVault.ID)

		// Return vault ID + snapshot + ALL peers + the wrapped DEK
		// New joiner gets everything in one request
		return c.JSON(fiber.Map{
			"vault_id":    targetVault.ID,
			"snapshot":    targetVault.Snapshot,
			"peers":       targetVault.Peers,
			"wrapped_dek": matched.WrappedDEK,
		})
	})

	// ── Message Routes (for offline delivery) ────────────────

	// POST /messages — Store message for offline peer
	app.Post("/messages", func(c *fiber.Ctx) error {
		var msg PendingMessage
		if err := c.BodyParser(&msg); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
		}

		msg.ID = generateID("msg_", 8)
		msg.CreatedAt = time.Now()
		msg.ExpiresAt = time.Now().Add(48 * time.Hour)

		store.mu.Lock()
		store.messages[msg.ForDeviceID] = append(
			store.messages[msg.ForDeviceID], &msg)
		store.mu.Unlock()

		return c.JSON(fiber.Map{"id": msg.ID, "success": true})
	})

	// GET /messages/:deviceID — Get pending messages
	app.Get("/messages/:deviceID", func(c *fiber.Ctx) error {
		store.mu.Lock()
		msgs := store.messages[c.Params("deviceID")]
		delete(store.messages, c.Params("deviceID"))
		store.mu.Unlock()

		if msgs == nil {
			msgs = []*PendingMessage{}
		}

		return c.JSON(fiber.Map{
			"messages": msgs,
			"count":    len(msgs),
		})
	})

	// DELETE /vaults/:id/peers/:deviceID — Remove peer from vault
	app.Delete("/vaults/:id/peers/:deviceID", func(c *fiber.Ctx) error {
		store.mu.Lock()
		defer store.mu.Unlock()

		vault, ok := store.vaults[c.Params("id")]
		if !ok {
			return c.Status(404).JSON(fiber.Map{
				"error": "vault not found",
			})
		}

		deviceID := c.Params("deviceID")
		found := false
		var remaining []Peer

		for _, p := range vault.Peers {
			if p.DeviceID == deviceID {
				found = true
				continue // skip — removes them
			}
			remaining = append(remaining, p)
		}

		if !found {
			return c.Status(404).JSON(fiber.Map{
				"error": "peer not found",
			})
		}

		vault.Peers = remaining
		vault.UpdatedAt = time.Now()

		log.Printf("🚫 Peer removed from vault %s: %s",
			c.Params("id"), deviceID)

		return c.JSON(fiber.Map{"success": true})
	})

	return app
}
