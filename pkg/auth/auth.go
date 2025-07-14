package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func NewToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func ValidateToken(token string) bool {
	// This function should validate the token and return the user ID if valid.
	// For now, we will return a placeholder user ID.
	return token == "valid_token_placeholder" // Invalid token returns 0 user ID
}
