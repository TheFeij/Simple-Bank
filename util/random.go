package util

import (
	"math/rand"
	"strings"
	"time"
)

var random *rand.Rand

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomInt(min int64, max int64) int64 {
	return random.Int63n(max-min+1) + min
}

func RandomString(length int) string {
	alphabet := "abcdefghijklmnopqrstuvwxyz"
	alphabetLength := len(alphabet)
	var randomString strings.Builder

	for i := 0; i < length; i++ {
		randomByte := alphabet[random.Intn(alphabetLength)]
		randomString.WriteByte(randomByte)
	}

	return randomString.String()
}

func RandomOwner() string {
	return RandomString(10)
}

//func RandomBalance() int64 {
//	return RandomInt(0, 10000)
//}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "IRL"}
	return currencies[random.Intn(len(currencies))]
}
