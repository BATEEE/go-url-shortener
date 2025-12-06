package utils

import (
	"crypto/rand"
	"math/big"
)

const alphabet = "9ABqrsCtuD4EFGHvwyxzIJKL5MNO6PQR0123SXTUVW78YZabcdefghijkmnopgl"

func GenerateRandomString(n int) string {
	b := make([]byte, n)
	max := big.NewInt(int64(len(alphabet)))

	for i := range b {
		num, _ := rand.Int(rand.Reader, max)
		b[i] = alphabet[num.Int64()]
	}
	return string(b)
}
