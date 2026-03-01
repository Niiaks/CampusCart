package tokenhash

import (
	"crypto/sha256"
	"encoding/hex"
)

// Hash returns a deterministic SHA-256 hex digest of the provided token.
func Hash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
