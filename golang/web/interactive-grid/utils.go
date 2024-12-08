package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func generateID() (string, error) {
	b := make([]byte, 10)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate secure random ID: %w", err)
	}

	return hex.EncodeToString(b), nil
}
