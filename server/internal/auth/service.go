package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"

	"github.com/zain-23/local-vault/server/internal/common/jwt"
	"github.com/zain-23/local-vault/server/internal/config"
)

// Argon2id params - controls how hard it is to brute-force password
const (
	argonTime		= 2				// iteration - more = slower but safer
	argonMemory		= 64 * 1024		// 64MB RAM - makes GPU attacks expensive
	argonThreads 	= 4				// parallel threads
	argonKeyLen		= 32			// output hash size in bytes
	argonSaltLen	= 16			// random salt size
)

// Service handles all auth business login
type Service struct {
	Store		*Store
	jwt 		*jwt.Service
	cfg			config.Config
}

func NewService(store *Store, jwtSvc *jwt.Service, cfg config.Config) *Service {
	return &Service{
		Store: store,
		jwt: jwtSvc,
		cfg: cfg,
	}
}


//  --------- Password hashing ---------
// HashPassword converts plain password to safe hash
func hashPssword(password string) (string, error) {
	salt := make([]byte, argonSaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// argon2.IDKey does the actual hashing - CPU + memory intensive by design
	hash := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	// Encode as text for storage - includes params so we can verify if we change setting later
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", 
			argon2.Version, argonMemory, argonTime, argonThreads, b64Salt, b64Hash), nil
}

// verifyPassword re-hashes with same salth and compares - returns true if password matches
func verifyPassword(password, endcode string) bool {
	parts := strings.Split(endcode, "$")
	if len(parts) != 6 {
		return  false
	}

	var memory uint32
	var iteration uint32
	var threads uint8
	fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iteration, &threads)

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return  false
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return  false
	}
	// Re-hash with same params + salt — should produce identical result
	computed := argon2.IDKey([]byte(password), salt, iteration, memory, threads, uint32(len(expectedHash)))
	

	// Constand-time compare - prevents timing attacks that measure comparison speed
	if len(computed) != len(expectedHash) {
		return  false
	}
	
	for i := range computed {
		if computed[i] != expectedHash[i] {
			return false
		}
	}
	return true
}

// --------------- Token helpers ---------------
// generateRandomToken create URL-safe random string - used for email links
func generateRandomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawStdEncoding.EncodeToString(b), nil
}

// sha25Hash hashes a token = we store the hash, never the raw token,
// so leaked DB can't reuse tokens
func sha25Hash(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:]) // h[:] converts fixed-size array to slice
}

// ------------- Signup ------------------------
// func (s *Service) Signup(ctx config.Config, req SignupRequest) (*SignupResponse, error) {
	
// } 