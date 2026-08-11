package helper

import (
	"crypto/rand"
	"encoding/hex"
)

func RandomHex(size int) string {
	b := make([]byte, size)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
