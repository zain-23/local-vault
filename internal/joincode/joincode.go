package joincode

import (
	"crypto/rand"
	"math/big"
	"strings"
)

// Alphabet matches server/device short codes (no 0/O/1/I/L).
const alphabet = "23456789ABCDEFGHJKMNPQRSTUVWXYZ"

// New returns a short "ABCD-1234" vault join code for email + WrapKey secret.
func New() (string, error) {
	const length = 8
	max := big.NewInt(int64(len(alphabet)))
	var sb strings.Builder
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		sb.WriteByte(alphabet[n.Int64()])
	}
	code := sb.String()
	return code[:4] + "-" + code[4:], nil
}
