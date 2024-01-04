package helper

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func GenerateDS() string {
	salt := "6s25p5ox5y14umn1p61aqyyvbvvl3lrt"
	t := time.Now().Unix()
	r := RandStringBytes(6)
	h := GetMD5Hash(fmt.Sprintf("salt=%s&t=%d&r=%s", salt, t, r))
	return fmt.Sprintf("%d,%s,%s", t, r, h)
}
