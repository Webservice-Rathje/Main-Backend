package utils

import (
	"crypto/sha512"
	"database/sql"
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

func CheckPasswordsMatch(currPassword string, conn *sql.DB, KundenID string) bool {
	stmt, _ := conn.Prepare("SELECT `Password` FROM `kunden` WHERE `KundenID`=?")
	type pwd_struct struct {
		Password string `json:"Password"`
	}
	resp, _ := stmt.Query(KundenID)
	var pwd pwd_struct
	for resp.Next() {
		err := resp.Scan(&pwd.Password)
		if err != nil {
			panic(err)
		}
	}
	salt := strings.Split(pwd.Password, "$")[0]
	currPasswordHash := HashPassword(currPassword, []byte(salt))
	resp.Close()
	stmt.Close()
	return pwd.Password == currPasswordHash
}
