package utils

import (
	"crypto/sha512"
	"encoding/base64"
	"math/rand"
	"strings"
	"time"
)

func GenerateSalt() []byte {
	rand.Seed(time.Now().Unix())
	charSet := "abcdedfghijklmnopqrstuvwxyzABCDEFHGIJKLMNOPQRSTUVWXYZ"
	var output strings.Builder
	length := 16
	for i := 0; i < length; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteString(string(randomChar))
	}
	return []byte(output.String())
}

func HashPassword(password string, salt []byte) string {
	var passwordBytes = []byte(password)
	var sha512Hasher = sha512.New()
	passwordBytes = append(passwordBytes, salt...)
	sha512Hasher.Write(passwordBytes)
	var hashedPasswordBytes = sha512Hasher.Sum(nil)
	var base64EncodedPasswordHash = base64.URLEncoding.EncodeToString(hashedPasswordBytes)
	return string(salt) + "$" + base64EncodedPasswordHash
}

func CheckPasswordsMatch(hashedPassword, currPassword string, salt []byte) bool {
	var currPasswordHash = HashPassword(currPassword, salt)

	return hashedPassword == currPasswordHash
}
