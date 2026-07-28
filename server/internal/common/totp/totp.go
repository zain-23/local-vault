package totp

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"github.com/pquerna/otp/totp"
)

const issuer = "LocalVault" // shown in the authenticator app

// GenerateSecret creates a new TOTP secret + an otpauth:// URL the frontend renders as a QR.
func GenerateSecret(accountEmail string) (secret, url string, err error) {
	key, err := totp.Generate(totp.GenerateOpts{Issuer: issuer, AccountName: accountEmail})
	if err != nil {
		return "", "", err
	}
	return key.Secret(), key.URL(), nil
}


// Validate reports whether code is a currently-valid 6-digit code for secret.
func Validate(secret, code string) bool {
	return totp.Validate(code, secret)
}

// GenerateBackupCodes returns n plaintext codes (shown once) and their hashes (stored).
func GenerateBackupCodes(n int) (plain, hashed []string, err error) {
	plain = make([]string, 0, n)
	hashed = make([]string, 0, n)
	for i := 0; i < n; i++ {
		b := make([]byte, 5) // 5 bytes → 10 hex chars
		if _, err := rand.Read(b); err != nil {
			return nil, nil, err
		}
		code := hex.EncodeToString(b)
		plain = append(plain, code)
		hashed = append(hashed, HashCode(code))
	}
	return plain, hashed, nil
}

// HashCode hashes a backup code — we store hashes, compare hashes, never the raw code.
func HashCode(code string) string {
	h := sha256.Sum256([]byte(code))
	return hex.EncodeToString(h[:])
}
