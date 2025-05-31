package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

var (
	ErrRequired = errors.New("empty inputed")
)

func DefaultString(s, fallback string) string {
	if s == "" {
		return fallback
	}

	return s
}

// hashToken encode token into hashed
func HashToken(token string) (string, error) {
	if token == "" {
		return "", ErrRequired
	}

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	return hashToken, nil
}
