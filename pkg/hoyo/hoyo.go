package hoyo

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"resin/pkg/config"
	"time"
)

func MakeDailyRequest(url string, ltoken string, ltuid string, actID string) (*http.Response, error) {
	jsonBody := []byte(fmt.Sprintf(`{"act_id": "%s"}`, actID))
	body := bytes.NewReader(jsonBody)
	r, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	r.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0")
	r.Header.Set("Accept", "application/json, text/plain, */*")
	r.Header.Set("Accept-Language", "en-US,en;q=0.5")
	r.Header.Set("Accept-Encoding", "gzip, deflate, br")
	r.Header.Set("Content-Type", "application/json;charset=utf-8")
	r.Header.Set("x-rpc-device_id", "7ba783da-1cc5-4c95-87f5-760e064faf37")
	r.Header.Set("x-rpc-app_version", "1.5.0")
	r.Header.Set("x-rpc-platform", "4")
	r.Header.Set("x-rpc-language", "en-us")
	r.Header.Set("x-rpc-device_name", "")
	r.Header.Set("Origin", "https://act.hoyolab.com")
	r.Header.Set("Connection", "keep-alive")
	r.Header.Set("Referer", "https://act.hoyolab.com/")
	r.Header.Set("Cookie", fmt.Sprintf("ltoken_v2=%s; ltuid_v2=%s", ltoken, ltuid))

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	response, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func MakeRequest(baseURL string, server string, genshinUID string, ltoken string, ltuid string) (*http.Response, error) {
	url := fmt.Sprintf("%s?server=%s&role_id=%s", baseURL, server, genshinUID)
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 11; SAMSUNG SM-G973U) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/14.2 Chrome/87.0.4280.141 Mobile Safari/537.36")
	r.Header.Set("Accept", "application/json, text/plain, */*")
	r.Header.Set("Accept-Language", "en-US,en;q=0.5")
	r.Header.Set("Accept-Encoding", "gzip, deflate, br")
	r.Header.Set("x-rpc-client_type", "5")
	r.Header.Set("x-rpc-app_version", "1.5.0")
	r.Header.Set("x-rpc-language", "en-us")
	r.Header.Set("Origin", "https://act.hoyolab.com")
	r.Header.Set("Connection", "keep-alive")
	r.Header.Set("Referer", "https://act.hoyolab.com/")
	r.Header.Set("Cookie", fmt.Sprintf("ltoken_v2=%s; ltuid_v2=%s", ltoken, ltuid))
	r.Header.Set("Sec-Fetch-Dest", "empty")
	r.Header.Set("Sec-Fetch-Mode", "cors")
	r.Header.Set("Sec-Fetch-Site", "same-site")

	client := &http.Client{}
	r.Header.Add("DS", GenerateDS())
	return client.Do(r)
}

func GetDailyData[T any](url string, ltoken string, ltuid string, actID string) (*T, error) {
	response, err := MakeDailyRequest(url, ltoken, ltuid, actID)
	if err != nil {
		return nil, err
	}
	return config.LoadJSON[T](response.Body)
}

func GetData[T any](baseURL string, server string, uid string, ltoken string, ltuid string) (*T, error) {
	response, err := MakeRequest(baseURL, server, uid, ltoken, ltuid)
	if err != nil {
		return nil, err
	}

	var reader io.ReadCloser

	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			return nil, err
		}
	default:
		reader = response.Body
	}

	defer reader.Close()

	return config.LoadJSON[T](reader)
}

func GenerateDS() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	random_bytes := make([]byte, 6)
	for i := range random_bytes {
		random_bytes[i] = charset[rand.Intn(len(charset))]
	}

	const salt = "6s25p5ox5y14umn1p61aqyyvbvvl3lrt"
	r := string(random_bytes)
	t := time.Now().Unix()
	text := fmt.Sprintf("salt=%s&t=%d&r=%s", salt, t, r)

	hash := md5.Sum([]byte(text))
	h := hex.EncodeToString(hash[:])

	return fmt.Sprintf("%d,%s,%s", t, r, h)
}

func GetTime(seconds int) (int, int) {
	hours := seconds / 3600
	minutes := (seconds / 60) - hours*60
	return hours, minutes
}
