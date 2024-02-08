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

func RandomPassword() string {
	password := RandomString(1, lowerCases)
	password += RandomString(1, upperCases)
	password += RandomString(1, numbers)
	password += RandomString(1, specials)
	password += RandomString(
		int(RandomInt(4, 60)),
		lowerCases+upperCases+numbers+numbers+specials)

	return password
}
func RandomBalance() int64 {
	return RandomInt(0, math.MaxInt64)
}
func RandomAmount() int32 {
	return int32(RandomInt(0, math.MinInt32))
}

func RandomFullname() string {
	randomString := RandomString(int(RandomInt(3, 64)), upperCases+lowerCases)

	index := RandomInt(1, int64(len(randomString)-2))
	fullname := randomString[:index]
	fullname += " "
	fullname += randomString[index+1:]

	return fullname
}

func RandomEmail() string {
	username := RandomUsername()
	domain := RandomString(int(RandomInt(5, 10)), lowerCases+upperCases)
	tld := RandomString(int(RandomInt(2, 5)), lowerCases)

	return username + "@" + domain + "." + tld
}
