package token

import (
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var randStringSrc = rand.NewSource(time.Now().UnixNano())

// RandomToken returns a (pseudo-)random token of length n
func RandomToken(n int) string {
	b := make([]byte, n)
	// A randStringSrc.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, randStringSrc.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randStringSrc.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}
