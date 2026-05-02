package main

// Signaling Server for LocalVault
//
// This is a completely separate program from the CLI.
// You deploy this once on a server.
// All LocalVault users share this one server.
//
// Think of it like a post office:
// - People drop off encrypted letters (invite codes, messages)
// - Recipients pick them up
// - Post office never opens the letters

import (
	// JSON encode/decode
	"fmt"
	"log" // server logging
	"os"
	"sync" // sync.RWMutex — safe concurrent map access
	"time"

	"github.com/gofiber/fiber/v2"                   // web framework like Express
	"github.com/gofiber/fiber/v2/middleware/cors"   // allow cross-origin requests
	"github.com/gofiber/fiber/v2/middleware/logger" // request logging
)

// ===== DATA STRUCTURES =====

// InviteSession represents a pending invite
// Created when Dev A runs: lv invite
// Consumed when Dev B runs: lv join LV-XXXX
type InviteSession struct {
	Code      string    `json:"code"`       // invite code e.g. "LV-A3F9-X2K1"
	DeviceID  string    `json:"device_id"`  // Dev A's device ID
	PublicKey []byte    `json:"public_key"` // Dev A's public key
	IPHint    string    `json:"ip_hint"`    // Dev A's IP (for LAN detection)
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"` // codes expire after 10 minutes
}

// PendingMessage represents an encrypted message
// stored when recipient is offline
// Like an email sitting in a mailbox
type PendingMessage struct {
	ID           string    `json:"id"`
	ForDeviceID  string    `json:"for_device_id"`  // recipient's device ID
	FromDeviceID string    `json:"from_device_id"` // sender's device ID
	Payload      []byte    `json:"payload"`        // encrypted blob (server cannot read)
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"` // messages expire after 48 hours
}

// PeerInfo is returned to Dev B after invite code lookup
// Contains everything B needs to connect to A
type PeerInfo struct {
	DeviceID  string `json:"device_id"`
	PublicKey []byte `json:"public_key"`
	IPHint    string `json:"ip_hint"`
}

// ===== IN-MEMORY STORE =====
// Simple storage using Go maps
// sync.RWMutex makes it safe for concurrent access
// RWMutex = multiple readers OR one writer at a time
// Like a read-write lock in any language

type Store struct {
	mu       sync.RWMutex                 // protects concurrent access
	invites  map[string]*InviteSession    // code → session
	messages map[string][]*PendingMessage // deviceID → messages
}

// NewStore creates empty store
func NewStore() *Store {
	return &Store{
		invites:  make(map[string]*InviteSession),
		messages: make(map[string][]*PendingMessage),
	}
}

// AddInvite stores a new invite session
// Called when Dev A runs: lv invite
func (s *Store) AddInvite(session *InviteSession) {
	// Lock for writing
	// Like: mutex.lock() in other languages
	s.mu.Lock()
	defer s.mu.Unlock() // unlock when function exits

	s.invites[session.Code] = session
}

// GetInvite retrieves and REMOVES an invite session
// Invite codes are single-use — once B joins, code is deleted
func (s *Store) GetInvite(code string) (*InviteSession, bool) {
	s.mu.Lock() // write lock because we delete after reading
	defer s.mu.Unlock()

	session, exists := s.invites[code]
	if !exists {
		return nil, false
	}

	// Check if expired
	if time.Now().After(session.ExpiresAt) {
		delete(s.invites, code) // clean up expired invite
		return nil, false
	}

	// Delete after retrieval — single use code
	delete(s.invites, code)
	return session, true
}

// AddMessage stores an encrypted message for offline peer
// Called when Dev A changes a secret and Dev B is offline
func (s *Store) AddMessage(msg *PendingMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Append to this device's message queue
	// Like pushing to an array: messages[deviceID].push(msg)
	s.messages[msg.ForDeviceID] = append(s.messages[msg.ForDeviceID], msg)
}

// GetMessages retrieves and clears all messages for a device
// Called when Dev B comes online: lv sync
func (s *Store) GetMessages(deviceID string) []*PendingMessage {
	s.mu.Lock()
	defer s.mu.Unlock()

	msgs := s.messages[deviceID]
	if msgs == nil {
		return []*PendingMessage{} // return empty slice not nil
	}

	// Clear messages after delivery
	// Like a mailbox — once you collect your mail, box is empty
	delete(s.messages, deviceID)
	return msgs
}

// Cleanup removes expired invites and messages
// Run periodically to keep memory usage low
func (s *Store) Cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()

	// Remove expired invites
	for code, session := range s.invites {
		if now.After(session.ExpiresAt) {
			delete(s.invites, code)
		}
	}

	// Remove expired messages
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

// ===== MAIN SERVER =====

func main() {
	// Create store
	store := NewStore()

	// Start cleanup goroutine
	// Runs every 5 minutes in background
	// Like: setInterval(() => cleanup(), 5 * 60 * 1000) in JS
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for range ticker.C { // range over channel = wait for each tick
			store.Cleanup()
			log.Println("🧹 Cleaned up expired sessions")
		}
	}()

	// Create Fiber app
	// Fiber is like Express.js but for Go
	app := fiber.New(fiber.Config{
		// Custom error handler
		// Returns JSON errors instead of HTML
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	// cors allows requests from any origin (CLI tool needs this)
	app.Use(cors.New())
	// logger prints each request to console
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${method} ${path} → ${status}\n",
	}))

	// ===== ROUTES =====

	// Health check — used to verify server is running
	// lv will ping this before any operation
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"version": "1.0.0",
			"time":    time.Now(),
		})
	})

	// POST /invite — Dev A creates an invite
	// Body: { code, device_id, public_key, ip_hint }
	// Response: { success: true }
	app.Post("/invite", func(c *fiber.Ctx) error {
		var session InviteSession

		// Parse JSON body
		// Like: const body = req.body in Express
		if err := c.BodyParser(&session); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "invalid request body",
			})
		}

		// Validate required fields
		if session.Code == "" || session.DeviceID == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "code and device_id are required",
			})
		}

		// Set timestamps
		session.CreatedAt = time.Now()
		session.ExpiresAt = time.Now().Add(10 * time.Minute) // expires in 10 min

		// Get real IP from request
		// c.IP() returns client IP address
		// Like: req.ip in Express
		if session.IPHint == "" {
			session.IPHint = c.IP()
		}

		store.AddInvite(&session)

		log.Printf("📨 Invite created: %s (device: %s)", session.Code, session.DeviceID)

		return c.JSON(fiber.Map{
			"success":    true,
			"expires_at": session.ExpiresAt,
		})
	})

	// GET /invite/:code — Dev B looks up an invite code
	// Response: { device_id, public_key, ip_hint }
	app.Get("/invite/:code", func(c *fiber.Ctx) error {
		// c.Params("code") gets URL parameter
		// Like: req.params.code in Express
		code := c.Params("code")

		session, exists := store.GetInvite(code)
		if !exists {
			return c.Status(404).JSON(fiber.Map{
				"error": "invite code not found or expired",
			})
		}

		log.Printf("🤝 Invite redeemed: %s", code)

		// Return peer info to Dev B
		return c.JSON(PeerInfo{
			DeviceID:  session.DeviceID,
			PublicKey: session.PublicKey,
			IPHint:    session.IPHint,
		})
	})

	// POST /messages — Dev A sends encrypted message for offline Dev B
	// Body: { for_device_id, from_device_id, payload }
	// Payload is encrypted — server cannot read it
	app.Post("/messages", func(c *fiber.Ctx) error {
		var msg PendingMessage

		if err := c.BodyParser(&msg); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "invalid request body",
			})
		}

		if msg.ForDeviceID == "" || msg.Payload == nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "for_device_id and payload are required",
			})
		}

		// Generate message ID
		msg.ID = fmt.Sprintf("msg_%d", time.Now().UnixNano())
		msg.CreatedAt = time.Now()
		msg.ExpiresAt = time.Now().Add(48 * time.Hour) // expires in 48 hours

		store.AddMessage(&msg)

		log.Printf("📬 Message queued for device: %s", msg.ForDeviceID)

		return c.JSON(fiber.Map{
			"success": true,
			"id":      msg.ID,
		})
	})

	// GET /messages/:deviceID — Dev B picks up pending messages
	// Called when Dev B comes online (lv sync)
	app.Get("/messages/:deviceID", func(c *fiber.Ctx) error {
		deviceID := c.Params("deviceID")

		msgs := store.GetMessages(deviceID)

		log.Printf("📮 Delivered %d messages to device: %s", len(msgs), deviceID)

		return c.JSON(fiber.Map{
			"messages": msgs,
			"count":    len(msgs),
		})
	})

	// GET /peers/:deviceID — check if a device is online
	// Used for real-time sync when both peers are online
	// For now returns basic info — will expand in Step 6 (daemon)
	app.Get("/peers/:deviceID", func(c *fiber.Ctx) error {
		// For now just return online status
		// We will expand this when we build the daemon
		return c.JSON(fiber.Map{
			"device_id": c.Params("deviceID"),
			"online":    false, // daemon will update this
		})
	})

	// Get port from environment variable
	// Like: process.env.PORT in Node.js
	// Default to 8080 if not set
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 LocalVault Signaling Server running on port %s", port)

	// Start server
	// Like: app.listen(port) in Express
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
