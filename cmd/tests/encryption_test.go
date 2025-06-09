package tests

import (
	"encoding/hex"
	"testing"

	"deployer.com/libs"
)

func TestEncryptionService_GenIv(t *testing.T) {
	es := libs.NewEncryptionService()
	iv := es.GenIv()
	if len(iv) != 24 { // 12 bytes in hex = 24 chars
		t.Errorf("Expected IV length 24, got %d", len(iv))
	}
	_, err := hex.DecodeString(iv)
	if err != nil {
		t.Errorf("IV is not valid hex: %v", err)
	}
}

func TestEncryptionService_EncryptDecrypt(t *testing.T) {
	es := libs.NewEncryptionService()
	key := make([]byte, 32) // AES-256
	for i := range key {
		key[i] = byte(i)
	}

	iv := es.GenIv()
	plain := "hello world"

	cipher, err := es.Encrypt(plain, iv)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	decrypted, err := es.Decrypt(cipher, iv)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}
	if decrypted != plain {
		t.Errorf("Expected decrypted text '%s', got '%s'", plain, decrypted)
	}
}

func TestEncryptionService_Decrypt_InvalidData(t *testing.T) {
	es := libs.NewEncryptionService()

	iv := es.GenIv()
	_, err := es.Decrypt("nothexdata", iv)
	if err == nil {
		t.Error("Expected error for invalid hex data, got nil")
	}
}

func TestEncryptionService_Encrypt_InvalidKey(t *testing.T) {
	es := libs.NewEncryptionService()
	iv := es.GenIv()
	_, err := es.Encrypt("data", iv)
	if err == nil {
		t.Error("Expected error for short key, got nil")
	}
}
