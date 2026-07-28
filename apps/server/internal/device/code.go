package device

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"math/big"
	"strings"
)

// userCodeAlphabet omits 0/O and 1/I/L — a human reads this code off a terminal
// and retypes it in a browser, and those are the characters they get wrong.
const userCodeAlphabet = "23456789ABCDEFGHJKMNPQRSTUVWXYZ"

// generateUserCode builds the short "ABCD-1234" code shown to the user.
// Guessing it only reveals an approval screen — it can never produce tokens.
func generateUserCode() (string, error) {
	const length = 8
	max := big.NewInt(int64(len(userCodeAlphabet))) // 31 — not a power of two

	var sb strings.Builder
	for i := 0; i < length; i++ {
		// rand.Int gives a uniform value in [0, max). Using rand.Read + modulo here
		// would bias toward early letters, because 256 % 31 != 0.
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		sb.WriteByte(userCodeAlphabet[n.Int64()])
	}

	code := sb.String()
	return code[:4] + "-" + code[4:], nil // hyphen purely for readability
}

// generateDeviceCode builds the long secret the CLI polls with. Independent of the
// user code by construction — a separate draw from crypto/rand, sharing no state.
func generateDeviceCode() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// URL-safe alphabet (-, _) so the code survives config files and headers intact.
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// sha256Hex hashes a code — we store the hash, never the raw code, so a leaked
// database yields nothing an attacker can poll with.
func sha256Hex(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:]) // h[:] converts the fixed-size array to a slice
}
