package vault

// This is the core of LocalVault
// Handles: storing secrets, reading secrets, encrypting vault file
// Think of this like your database layer

import (
	"encoding/json" // like JSON.parse / JSON.stringify in JS
	"errors"
	"fmt"
	"os"            // file system operations
	"path/filepath" // cross-platform file paths (handles Windows \ vs Unix /)
	"time"

	"github.com/zain-23/local-vault/internal/crypto"
)

// ===== DATA STRUCTURES =====
// In Go, structs are like JS objects with fixed shape (like TypeScript interfaces)

// Secret represents one environment variable
// The `json:"..."` tags tell Go how to name fields when converting to/from JSON
// Same as: type Secret = { key: string, value: string, env: string }
type Secret struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	Env       string    `json:"env"`        // "development", "staging", "production"
	UpdatedAt time.Time `json:"updated_at"` // when was this last changed
	UpdatedBy string    `json:"updated_by"` // which device changed it (for sync)
}

// VaultFile is what gets written to disk (encrypted)
// Contains all secrets + metadata
type VaultFile struct {
	Version   string    `json:"version"`   // vault format version
	Salt      []byte    `json:"salt"`      // salt for key derivation
	PassHash  string    `json:"pass_hash"` // hashed passphrase for validation
	Secrets   []Secret  `json:"secrets"`   // all the secrets
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Vault is the in-memory representation
// After we decrypt + parse VaultFile, we work with this struct
// It holds both the data and the paths/keys needed to save it back
type Vault struct {
	file       VaultFile // the actual data
	passphrase string    // kept in memory to re-encrypt on save
	dir        string    // path to .lv/ folder
}

// ===== CONSTANTS =====

const (
	lvDir     = ".lv"            // hidden folder name
	vaultFile = "vault.json.enc" // encrypted vault filename
	version   = "1.0.0"
)

// ===== PUBLIC FUNCTIONS =====

// Init creates a new vault in the given directory
// Called by: lv init
// JS equivalent: async function initVault(dir, passphrase)
func Init(dir string, passphrase string) error {
	// Build path to .lv folder
	// filepath.Join handles slashes correctly on all OS
	lvPath := filepath.Join(dir, lvDir)

	// Check if vault already exists
	// os.Stat returns file info or error if not found
	if _, err := os.Stat(lvPath); err == nil {
		return errors.New("vault already exists in this directory")
	}

	// Create .lv directory
	// 0700 = only owner can read/write/execute (Unix permissions)
	// Like: chmod 700 .lv
	if err := os.MkdirAll(lvPath, 0700); err != nil {
		return fmt.Errorf("failed to create vault directory: %w", err)
	}

	// Generate random salt for key derivation
	salt, err := crypto.GenerateSalt()
	if err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	// Create empty vault file structure
	vf := VaultFile{
		Version:   version,
		Salt:      salt,
		PassHash:  crypto.HashPassphrase(passphrase),
		Secrets:   []Secret{}, // empty slice (like [] in JS)
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create vault instance
	v := &Vault{ // & means we're creating a pointer to Vault
		file:       vf,
		passphrase: passphrase,
		dir:        lvPath,
	}

	// Save encrypted vault to disk
	if err := v.save(); err != nil {
		return fmt.Errorf("failed to save vault: %w", err)
	}

	// Add .gitignore rules to prevent secrets from being committed
	if err := addGitignoreRules(dir); err != nil {
		// Non-fatal — vault still works, just warn user
		fmt.Println("⚠️  Could not update .gitignore:", err)
	}

	return nil
}

// Load reads and decrypts an existing vault from disk
// Called before every command that needs to read/write secrets
// JS equivalent: async function loadVault(dir, passphrase)
func Load(dir string, passphrase string) (*Vault, error) {
	lvPath := filepath.Join(dir, lvDir)
	vaultPath := filepath.Join(lvPath, vaultFile)

	// Read encrypted file from disk
	encryptedData, err := os.ReadFile(vaultPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("no vault found. run 'lv init' first")
		}
		return nil, fmt.Errorf("failed to read vault: %w", err)
	}

	// Derive encryption key from passphrase
	// We need the salt first — it's stored unencrypted at start of file
	// First 16 bytes are salt (we'll improve this later with a header)
	// For now, decrypt with a temporary approach:
	// Actually we store salt in the JSON so we need to do this in 2 steps

	// Step 1: Try to decrypt with provided passphrase
	// We need the salt — temporarily store it unencrypted as first 16 bytes
	if len(encryptedData) < 16 {
		return nil, errors.New("vault file is corrupted")
	}

	salt := encryptedData[:16]
	actualEncrypted := encryptedData[16:]

	// Derive the key using the stored salt
	key := crypto.DeriveKey(passphrase, salt)

	// Decrypt the vault data
	decryptedData, err := crypto.Decrypt(actualEncrypted, key)
	if err != nil {
		return nil, errors.New("wrong passphrase or corrupted vault")
	}

	// Parse JSON into our VaultFile struct
	// Like: const vaultFile = JSON.parse(decryptedData)
	var vf VaultFile
	if err := json.Unmarshal(decryptedData, &vf); err != nil {
		return nil, fmt.Errorf("failed to parse vault: %w", err)
	}

	// Verify passphrase matches stored hash
	if crypto.HashPassphrase(passphrase) != vf.PassHash {
		return nil, errors.New("wrong passphrase")
	}

	return &Vault{
		file:       vf,
		passphrase: passphrase,
		dir:        lvPath,
	}, nil
}

// Add stores a new secret in the vault
// Called by: lv add KEY=value
func (v *Vault) Add(key, value, env string) error {
	// Check if key already exists — update it if so
	for i, s := range v.file.Secrets {
		if s.Key == key && s.Env == env {
			// Update existing secret
			v.file.Secrets[i].Value = value
			v.file.Secrets[i].UpdatedAt = time.Now()
			return v.save()
		}
	}

	// Add new secret
	v.file.Secrets = append(v.file.Secrets, Secret{
		Key:       key,
		Value:     value,
		Env:       env,
		UpdatedAt: time.Now(),
	})

	v.file.UpdatedAt = time.Now()
	return v.save()
}

// Get retrieves a single secret by key
// Called by: lv get KEY
func (v *Vault) Get(key, env string) (string, error) {
	for _, s := range v.file.Secrets {
		if s.Key == key && (s.Env == env || env == "") {
			return s.Value, nil
		}
	}
	return "", fmt.Errorf("secret '%s' not found", key)
}

// List returns all secrets for a given environment
// Called by: lv list
func (v *Vault) List(env string) []Secret {
	if env == "" {
		return v.file.Secrets // return all
	}

	// Filter by environment
	// In JS: secrets.filter(s => s.env === env)
	var filtered []Secret
	for _, s := range v.file.Secrets {
		if s.Env == env {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

// Remove deletes a secret from the vault
// Called by: lv remove KEY
func (v *Vault) Remove(key, env string) error {
	found := false
	// Build new slice without the removed secret
	// In JS: secrets.filter(s => !(s.key === key && s.env === env))
	var remaining []Secret
	for _, s := range v.file.Secrets {
		if s.Key == key && (s.Env == env || env == "") {
			found = true
			continue // skip this one (removes it)
		}
		remaining = append(remaining, s)
	}

	if !found {
		return fmt.Errorf("secret '%s' not found", key)
	}

	v.file.Secrets = remaining
	v.file.UpdatedAt = time.Now()
	return v.save()
}

// InjectMap returns secrets as a simple key→value map
// Used by lv inject to export to shell or child process
func (v *Vault) InjectMap(env string) map[string]string {
	result := make(map[string]string) // like {} in JS
	for _, s := range v.file.Secrets {
		if s.Env == env || s.Env == "" || env == "" {
			result[s.Key] = s.Value
		}
	}
	return result
}

// ImportEnvFile reads a .env file and adds all its secrets to vault
// Called by: lv import .env.local
func (v *Vault) ImportEnvFile(path, env string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("failed to read file: %w", err)
	}

	count := 0
	lines := splitLines(string(data))

	for _, line := range lines {
		// Skip comments and empty lines
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		// Parse KEY=VALUE format
		key, value, found := parseEnvLine(line)
		if !found {
			continue
		}

		if err := v.Add(key, value, env); err != nil {
			return count, err
		}
		count++
	}

	return count, nil
}

// ===== PRIVATE FUNCTIONS =====
// In Go, lowercase function names are private to the package
// Like unexported functions in a JS module

// save encrypts and writes vault to disk
func (v *Vault) save() error {
	// Convert struct to JSON bytes
	// Like: JSON.stringify(vaultFile)
	jsonData, err := json.Marshal(v.file)
	if err != nil {
		return fmt.Errorf("failed to serialize vault: %w", err)
	}

	// Derive encryption key from passphrase + salt
	key := crypto.DeriveKey(v.passphrase, v.file.Salt)

	// Encrypt the JSON data
	encrypted, err := crypto.Encrypt(jsonData, key)
	if err != nil {
		return fmt.Errorf("failed to encrypt vault: %w", err)
	}

	// Prepend salt to encrypted data
	// Format: [16 bytes salt][encrypted data]
	// We need salt to derive key during Load()
	finalData := append(v.file.Salt, encrypted...)

	// Write to disk with restrictive permissions
	// 0600 = only owner can read/write
	vaultPath := filepath.Join(v.dir, vaultFile)
	return os.WriteFile(vaultPath, finalData, 0600)
}

// addGitignoreRules adds .lv/ entries to .gitignore
func addGitignoreRules(dir string) error {
	gitignorePath := filepath.Join(dir, ".gitignore")

	rules := "\n# LocalVault\n.lv/vault.json.enc\n.lv/identity.key\n"

	// Open file in append mode, create if not exists
	// os.O_APPEND = add to end, os.O_CREATE = create if missing
	// os.O_WRONLY = write only
	f, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close() // defer = runs when function exits (like finally in JS)

	_, err = f.WriteString(rules)
	return err
}

// splitLines splits text into lines (handles \r\n on Windows)
func splitLines(text string) []string {
	var lines []string
	start := 0
	for i, c := range text {
		if c == '\n' {
			line := text[start:i]
			// Trim \r for Windows line endings
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			lines = append(lines, line)
			start = i + 1
		}
	}
	if start < len(text) {
		lines = append(lines, text[start:])
	}
	return lines
}

// parseEnvLine parses "KEY=VALUE" into separate key and value
// Returns false if line doesn't match expected format
func parseEnvLine(line string) (string, string, bool) {
	for i, c := range line {
		if c == '=' {
			key := line[:i]
			value := line[i+1:]
			// Strip surrounding quotes from value if present
			// Handles: KEY="value" or KEY='value'
			if len(value) >= 2 {
				if (value[0] == '"' && value[len(value)-1] == '"') ||
					(value[0] == '\'' && value[len(value)-1] == '\'') {
					value = value[1 : len(value)-1]
				}
			}
			return key, value, true
		}
	}
	return "", "", false
}
