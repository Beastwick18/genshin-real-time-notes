package helper

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
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

func GenerateRequest(url string, ltoken string, ltuid string) *http.Request {
	rOrig, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	rOrig.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 11; SAMSUNG SM-G973U) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/14.2 Chrome/87.0.4280.141 Mobile Safari/537.36")
	rOrig.Header.Add("Accept", "application/json, text/plain, */*")
	rOrig.Header.Add("Accept-Language", "en-US,en;q=0.5")
	rOrig.Header.Add("Accept-Encoding", "gzip, deflate, br")
	rOrig.Header.Add("x-rpc-client_type", "5")
	rOrig.Header.Add("x-rpc-app_version", "1.5.0")
	rOrig.Header.Add("x-rpc-language", "en-us")
	rOrig.Header.Add("Origin", "https://act.hoyolab.com")
	rOrig.Header.Add("Connection", "keep-alive")
	rOrig.Header.Add("Referer", "https://act.hoyolab.com/")
	rOrig.Header.Add("Cookie", fmt.Sprintf("ltoken_v2=%s; ltuid_v2=%s", ltoken, ltuid))
	rOrig.Header.Add("Sec-Fetch-Dest", "empty")
	rOrig.Header.Add("Sec-Fetch-Mode", "cors")
	rOrig.Header.Add("Sec-Fetch-Site", "same-site")
	return rOrig
}
