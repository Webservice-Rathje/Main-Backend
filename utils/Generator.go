package utils

import (
	"math/rand"
	"strings"
	"time"
)

func ID_generator(length int) string {
	rand.Seed(time.Now().Unix())
	charSet := "0123456789"
	var output strings.Builder
	for i := 0; i < length; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteString(string(randomChar))
	}
	return output.String()
}

func Generate_KundenID() string {
	conn := GetConn()
	stmt, _ := conn.Prepare("SELECT * FROM `kunden` WHERE `KundenID`=?;")
	var kid string
	for {
		kid = ID_generator(7)
		resp, _ := stmt.Query(kid)
		exists := false
		for resp.Next() {
			exists = true
		}
		resp.Close()
		if exists == false {
			break
		}
	}
	stmt.Close()
	conn.Close()
	return kid
}

func Generate_2FA_Code() string {
	conn := GetConn()
	stmt, _ := conn.Prepare("SELECT * FROM `2FA-Codes` WHERE `Code`=?;")
	var code string
	for {
		code = ID_generator(6)
		resp, _ := stmt.Query(code)
		exists := false
		for resp.Next() {
			exists = true
		}
		resp.Close()
		if exists == false {
			break
		}
	}
	stmt.Close()
	conn.Close()
	return code
}

func GenerateToken() string {
	rand.Seed(time.Now().Unix())
	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var output strings.Builder
	for i := 0; i < 64; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteString(string(randomChar))
	}
	return output.String()
}
