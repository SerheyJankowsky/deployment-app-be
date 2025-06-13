package libs

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateApiKey() string {
	apiKey := make([]byte, 32)
	rand.Read(apiKey)
	return hex.EncodeToString(apiKey)
}
