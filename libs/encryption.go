package libs

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

type EncryptionService struct {
}

func NewEncryptionService() *EncryptionService {
	return &EncryptionService{}
}

func (e *EncryptionService) GenIv() string {
	iv := make([]byte, 12) // GCM standard nonce size
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	return hex.EncodeToString(iv)
}

func (e *EncryptionService) Encrypt(data string, iv string) (string, error) {
	key, err := hex.DecodeString(os.Getenv("ENCRYPTION_KEY"))
	if err != nil {
		return "", err
	}
	if len(key) != 32 {
		return "", fmt.Errorf("invalid key length: got %d, want 32", len(key))
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	ivBytes, err := hex.DecodeString(iv)
	if err != nil {
		return "", err
	}
	if len(ivBytes) != 12 {
		return "", fmt.Errorf("invalid IV length: got %d, want 12", len(ivBytes))
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	cipherText := gcm.Seal(nil, ivBytes, []byte(data), nil)
	return hex.EncodeToString(cipherText), nil
}

func (e *EncryptionService) Decrypt(data string, iv string) (string, error) {
	if data == "" || data == "null" {
		return "", nil
	}
	key, err := hex.DecodeString(os.Getenv("ENCRYPTION_KEY"))
	if err != nil {
		return "", err
	}
	if len(key) != 32 {
		return "", fmt.Errorf("invalid key length: got %d, want 32", len(key))
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	ivBytes, err := hex.DecodeString(iv)
	if err != nil {
		return "", err
	}
	if len(ivBytes) != 12 {
		return "", fmt.Errorf("invalid IV length: got %d, want 12", len(ivBytes))
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	cipherText, err := hex.DecodeString(data)
	if err != nil {
		return "", err
	}
	plaintext, err := gcm.Open(nil, ivBytes, cipherText, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
