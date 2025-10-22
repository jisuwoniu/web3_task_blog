package utils

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestGenerateRandomBytes(t *testing.T) {
	length := 16
	bytes, err := GenerateRandomBytes(length)
	if err != nil {
		t.Fatalf("GenerateRandomBytes failed: %v", err)
	}

	if len(bytes) != length {
		t.Fatalf("Expected %d bytes, got %d", length, len(bytes))
	}

	// Test error case
	_, err = GenerateRandomBytes(0)
	if err == nil {
		t.Fatal("Expected error for length <= 0, got nil")
	}
}

func TestGenerateBase64Key(t *testing.T) {
	length := 64
	key, err := GenerateBase64Key(length)
	if err != nil {
		t.Fatalf("GenerateBase64Key failed: %v", err)
	}

	fmt.Println(key)
	decoded, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		t.Fatalf("Failed to decode Base64 key: %v", err)
	}
	fmt.Println(decoded)
	if len(decoded) != length {
		t.Fatalf("Expected %d bytes, got %d", length, len(decoded))
	}

	// Test error case
	_, err = GenerateBase64Key(0)
	if err == nil {
		t.Fatal("Expected error for length <= 0, got nil")
	}
}

func TestGenerateHexKey(t *testing.T) {
	length := 16
	key, err := GenerateHexKey(length)
	if err != nil {
		t.Fatalf("GenerateHexKey failed: %v", err)
	}

	decoded, err := hex.DecodeString(key)
	if err != nil {
		t.Fatalf("Failed to decode Hex key: %v", err)
	}

	if len(decoded) != length {
		t.Fatalf("Expected %d bytes, got %d", length, len(decoded))
	}

	// Test error case
	_, err = GenerateHexKey(0)
	if err == nil {
		t.Fatal("Expected error for length <= 0, got nil")
	}
}
