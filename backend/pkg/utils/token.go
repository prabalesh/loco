package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateToken(outputLength int) (string, error) {
	if outputLength <= 0 {
		return "", fmt.Errorf("invalid token length")
	}

	byteLen := (outputLength * 3) / 4
	if byteLen == 0 {
		byteLen = 1
	}

	b := make([]byte, byteLen)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	token := base64.RawURLEncoding.EncodeToString(b)

	if len(token) > outputLength {
		token = token[:outputLength]
	}

	return token, nil
}
