package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// RandomHex generates a cryptographically secure random hex string
// with a length of 2*n characters (because n bytes = 2n hex chars).
//
// Example:
//
//	RandomHex(32) → 64-character secure token
//
// Common uses:
//   - Refresh tokens
//   - Password reset tokens
//   - API keys
//
// Returns an error if the system's secure random generator fails.
func RandomHex(byteLength int) (string, error) {
	randomBytes := make([]byte, byteLength)

	// rand.Read fills the slice with cryptographically secure random bytes.
	// This is safe for generating authentication tokens.
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Convert bytes to hexadecimal string.
	return hex.EncodeToString(randomBytes), nil
}

// SHA256Hex calculates the SHA-256 hash of a string and returns it as hex.
//
// This is useful for:
//   - Hashing refresh tokens so the DB never stores raw tokens
//   - Generating consistent identifiers
//
// Note: This is not used for passwords. Passwords must use bcrypt.
func SHA256Hex(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}
