package utils

import (
	"crypto/rand"
	"encoding/base64"
	"os"
)

const (
	tokenLength int = 150
)

// GenerateRandomToken generates a random token of the specified length.
func GenerateRandomToken() (string, []byte, error) {
	// Salt string used to encode the token used to verify the user
	saltString := os.Getenv("EMAIL_VERIFICATION_SALT")

	// Generate a random token
	tokenBytes := make([]byte, tokenLength)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", nil, err
	}

	// Combine the token with the salt and encode using base64
	tokenWithSalt := base64.URLEncoding.EncodeToString([]byte(saltString + base64.URLEncoding.EncodeToString(tokenBytes)))

	return tokenWithSalt, tokenBytes, nil
}
