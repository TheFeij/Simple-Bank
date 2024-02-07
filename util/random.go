package util

import (
	"math"
	"math/rand"
	"strings"
	"time"
)

var lowerCases = "abcdefghijklmnopqrstuvwxyz"
var upperCases = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
var numbers = "0123456789"
var specials = "_!@#$%&*^"

var random *rand.Rand

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomInt(min int64, max int64) int64 {
	return random.Int63n(max-min+1) + min
}

func RandomString(length int, alphabet string) string {
	alphabetLength := len(alphabet)
	var randomString strings.Builder

	for i := 0; i < length; i++ {
		randomByte := alphabet[random.Intn(alphabetLength)]
		randomString.WriteByte(randomByte)
	}

	return randomString.String()
}

func RandomID() int64 {
	return RandomInt(1, math.MaxInt64)
}

func RandomUsername() string {
	username := RandomString(1, lowerCases+upperCases)
	username += RandomString(
		int(RandomInt(3, 62)),
		lowerCases+upperCases+numbers+"_")
	username += RandomString(
		int(RandomInt(0, 1)),
		lowerCases+upperCases+numbers)

	return username
}
