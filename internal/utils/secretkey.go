package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
)

// GenerateRandomBytes generates a slice of random bytes of the specified length.
func GenerateRandomBytes(length int) ([]byte, error) {
	if length <= 0 {
		return nil, errors.New("length must be greater than 0")
	}

	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}

	return bytes, nil
}

// GenerateBase64Key generates a Base64-encoded random key of the specified length.
func GenerateBase64Key(length int) (string, error) {
	bytes, err := GenerateRandomBytes(length)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bytes), nil
}

// GenerateHexKey generates a hex-encoded random key of the specified length.
func GenerateHexKey(length int) (string, error) {
	bytes, err := GenerateRandomBytes(length)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}