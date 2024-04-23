package internal

import "time"
import mathrand "math/rand"

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"0123456789"

var seededRand = mathrand.New(mathrand.NewSource(time.Now().UnixNano()))

func stringWithCharset(length int, charset string) []byte {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return b
}

func RandomString(length int) []byte {
	return stringWithCharset(length, charset)
}
