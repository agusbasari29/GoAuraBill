package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// GenerateSignature membuat signature HMAC-SHA256 untuk Tripay
func GenerateSignature(merchantCode, merchantRef string, amount int, privateKey string) string {
	data := merchantCode + merchantRef + string(rune(amount))
	mac := hmac.New(sha256.New, []byte(privateKey))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}