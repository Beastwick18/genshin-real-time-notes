package hoyo

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"resin/pkg/config"
	"resin/pkg/helper"
	"time"
)

func MakeRequest(baseURL string, server string, genshinUID string, ltoken string, ltuid string) *http.Response {
	url := fmt.Sprintf("%s?server=%s&role_id=%s", baseURL, server, genshinUID)
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
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
	return response
}

func GetData[T any](baseURL string, server string, uid string, ltoken string, ltuid string) T {
	response := MakeRequest(baseURL, server, uid, ltoken, ltuid)

	// Check that the server actually sent compressed data
	var reader io.ReadCloser
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(response.Body)
	default:
		reader = response.Body
	}
	defer reader.Close()

	return config.LoadJSON[T](reader)
}

func GenerateDS() string {
	salt := "6s25p5ox5y14umn1p61aqyyvbvvl3lrt"
	t := time.Now().Unix()
	r := helper.RandStringBytes(6)
	h := helper.GetMD5Hash(fmt.Sprintf("salt=%s&t=%d&r=%s", salt, t, r))
	return fmt.Sprintf("%d,%s,%s", t, r, h)
}

func GetTime(seconds int) (int, int) {
	hours := seconds / 3600
	minutes := (seconds / 60) - hours*60
	return hours, minutes
}
