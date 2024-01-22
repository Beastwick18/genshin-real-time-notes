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
	r.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0")
	r.Header.Add("Accept", "application/json, text/plain, */*")
	r.Header.Add("Accept-Language", "en-US,en;q=0.5")
	r.Header.Add("Accept-Encoding", "gzip, deflate, br")
	r.Header.Add("Content-Type", "application/json;charset=utf-8")
	r.Header.Add("x-rpc-device_id", "7ba783da-1cc5-4c95-87f5-760e064faf37")
	r.Header.Add("x-rpc-app_version", "1.5.0")
	r.Header.Add("x-rpc-platform", "4")
	r.Header.Add("x-rpc-language", "en-us")
	r.Header.Add("x-rpc-device_name", "")
	r.Header.Add("Origin", "https://act.hoyolab.com")
	r.Header.Add("Connection", "keep-alive")
	r.Header.Add("Referer", "https://act.hoyolab.com/")
	r.Header.Add("Cookie", fmt.Sprintf("ltoken_v2=%s; ltuid_v2=%s", ltoken, ltuid))

	client := &http.Client{}
	response, err := client.Do(r)
	return response, nil
}

func MakeRequest(baseURL string, server string, genshinUID string, ltoken string, ltuid string) (*http.Response, error) {
	url := fmt.Sprintf("%s?server=%s&role_id=%s", baseURL, server, genshinUID)
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 11; SAMSUNG SM-G973U) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/14.2 Chrome/87.0.4280.141 Mobile Safari/537.36")
	r.Header.Add("Accept", "application/json, text/plain, */*")
	r.Header.Add("Accept-Language", "en-US,en;q=0.5")
	r.Header.Add("Accept-Encoding", "gzip, deflate, br")
	r.Header.Add("x-rpc-client_type", "5")
	r.Header.Add("x-rpc-app_version", "1.5.0")
	r.Header.Add("x-rpc-language", "en-us")
	r.Header.Add("Origin", "https://act.hoyolab.com")
	r.Header.Add("Connection", "keep-alive")
	r.Header.Add("Referer", "https://act.hoyolab.com/")
	r.Header.Add("Cookie", fmt.Sprintf("ltoken_v2=%s; ltuid_v2=%s", ltoken, ltuid))
	r.Header.Add("Sec-Fetch-Dest", "empty")
	r.Header.Add("Sec-Fetch-Mode", "cors")
	r.Header.Add("Sec-Fetch-Site", "same-site")

	client := &http.Client{}
	r.Header.Add("DS", GenerateDS())
	response, err := client.Do(r)
	return response, nil
}

func GetDailyData[T any](url string, ltoken string, ltuid string, actID string) (*T, error) {
	response, err := MakeDailyRequest(url, ltoken, ltuid, actID)
	if err != nil {
		return nil, err
	}
	json, err := config.LoadJSON[T](response.Body)
	if err != nil {
		return nil, err
	}
	return json, nil
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

	json, err := config.LoadJSON[T](reader)
	if err != nil {
		return nil, err
	}
	return json, nil
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
