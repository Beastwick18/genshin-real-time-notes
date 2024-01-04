package main

import (
	// "bytes"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"genshin/config"
	"genshin/helper"
	"genshin/icon"
	"io"
	"net/http"
	"time"

	"github.com/getlantern/systray"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("Genshin Real-Time Notes")
	systray.SetTooltip("?/?")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	cfg := config.LoadConfig("./config.json")

	go func() {
		posturl := fmt.Sprintf("https://bbs-api-os.hoyolab.com/game_record/genshin/api/dailyNote?server=%s&role_id=%s", cfg.Server, cfg.Genshin_uuid)

		// Create a HTTP post request
		rOrig, err := http.NewRequest("GET", posturl, nil)
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
		rOrig.Header.Add("Cookie", fmt.Sprintf("ltoken_v2=%s; ltuid_v2=%s", cfg.Ltoken, cfg.Ltuid))
		rOrig.Header.Add("Sec-Fetch-Dest", "empty")
		rOrig.Header.Add("Sec-Fetch-Mode", "cors")
		rOrig.Header.Add("Sec-Fetch-Site", "same-site")

		client := &http.Client{}
		for {
			func() {
				r := rOrig
				r.Header.Add("DS", helper.GenerateDS())
				response, err := client.Do(r)
				if err != nil {
					panic(err)
				}
				defer response.Body.Close()

				// Check that the server actually sent compressed data
				var reader io.ReadCloser
				switch response.Header.Get("Content-Encoding") {
				case "gzip":
					reader, err = gzip.NewReader(response.Body)
					defer reader.Close()
				default:
					reader = response.Body
				}

				jsonMap := make(map[string]interface{})
				buf := new(bytes.Buffer)
				buf.ReadFrom(reader)
				jsonString := buf.String()
				err = json.Unmarshal([]byte(jsonString), &jsonMap)
				if err != nil {
					panic(err)
				}
				fmt.Println(jsonMap)

				current := "?"
				max := "?"
				full := false
				if data, ok := jsonMap["data"].(map[string]interface{}); ok {
					if resin, ok := data["current_resin"]; ok {
						current = fmt.Sprint(resin)
					}
					if resin, ok := data["max_resin"]; ok {
						max = fmt.Sprint(resin)
						full = current == max
					}
				}
				if full {
					systray.SetIcon(icon.FullData)
				} else {
					systray.SetIcon(icon.NotFullData)
				}
				systray.SetTooltip(fmt.Sprintf("%s/%s", current, max))
				time.Sleep(10 * time.Second)
			}()
		}
	}()

	go func() {
		mQuit.SetIcon(icon.Data)
		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
				fmt.Println("quitting")
				return
			}
		}
	}()
}

func onExit() {
	// clean up here
}
