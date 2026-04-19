package util

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRandomCode() (string, error) {
	b := make([]byte, 3)
	_, err := rand.Read(b)
	return hex.EncodeToString(b), err
}
