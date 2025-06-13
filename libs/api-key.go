package libs

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

func GenerateApiKey() string {
	// Generate 32 bytes (256 bits) of random data
	apiKey := make([]byte, 32)
	if _, err := rand.Read(apiKey); err != nil {
		panic(fmt.Sprintf("Failed to generate API key: %v", err))
	}
	return hex.EncodeToString(apiKey)
}

// GenerateSecureApiKey generates a more secure API key with timestamp prefix
func GenerateSecureApiKey() string {
	// Add timestamp prefix for uniqueness and tracking
	timestamp := time.Now().Unix()

	// Generate 28 bytes of random data (to keep total length reasonable)
	randomBytes := make([]byte, 28)
	if _, err := rand.Read(randomBytes); err != nil {
		panic(fmt.Sprintf("Failed to generate secure API key: %v", err))
	}

	// Combine timestamp and random bytes
	return fmt.Sprintf("ak_%d_%s", timestamp, hex.EncodeToString(randomBytes))
}
