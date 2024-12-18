package util

import (
	"crypto/rand"
	"encoding/base64"
)

func RandomToken(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
