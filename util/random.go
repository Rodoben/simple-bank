package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"
const digits = "0123456789"
const specialChars = `!@#$%^&*()_+-={}:<>?/.,;`
const ALPHABET = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	rand.New(rand.NewSource(time.Now().Unix()))
}

func RandomInt(min, max int64) int64 {
	return rand.Int63n(max-min) + min
}

func RandomString(n int) string {
	var sb strings.Builder
	a := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(a)]
		sb.WriteByte(c)
	}
	return sb.String()

}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(int64(1), int64(10))
}

func RandomCurrency() string {
	currencies := []string{"INR", "USD", "AED", "TBH"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
