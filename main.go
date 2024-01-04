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
	"strconv"
	"time"

	"github.com/getlantern/systray"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(icon.NotFullData)
	systray.SetTitle("Genshin Real-Time Notes")
	systray.SetTooltip("?/?")
	mResin := systray.AddMenuItem("Resin: ?/?", "Resin")
	mResin.SetIcon(icon.NotFullData)

	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	cfg := config.LoadConfig("./config.json")

	go func() {
		url := fmt.Sprintf("https://bbs-api-os.hoyolab.com/game_record/genshin/api/dailyNote?server=%s&role_id=%s", cfg.Server, cfg.Genshin_uuid)

		// Create a HTTP request
		rOrig := helper.GenerateRequest(url, cfg.Ltoken, cfg.Ltuid)

		client := &http.Client{}
		for {
			func() {
				r := rOrig
				r.Header.Add("DS", helper.GenerateDS())
				response, err := client.Do(r)
				if err != nil {
					panic(err)
				}

				// Check that the server actually sent compressed data
				var reader io.ReadCloser
				switch response.Header.Get("Content-Encoding") {
				case "gzip":
					reader, err = gzip.NewReader(response.Body)
				default:
					reader = response.Body
				}
				defer reader.Close()

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
				recovery := ""
				full := false
				if data, ok := jsonMap["data"].(map[string]interface{}); ok {
					if resin, ok := data["current_resin"]; ok {
						current = fmt.Sprint(resin)
					}
					if resin, ok := data["max_resin"]; ok {
						max = fmt.Sprint(resin)
						full = current == max
					}
					if rec_time, ok := data["resin_recovery_time"]; ok {
						seconds, _ := strconv.Atoi(fmt.Sprint(rec_time))
						if seconds != 0 {
							hours := seconds / 3600
							minutes := (seconds / 60) - hours*60
							recovery = fmt.Sprintf(" [%dh %dm]", hours, minutes)
						}
					}
				}
				if full {
					systray.SetIcon(icon.FullData)
				} else {
					systray.SetIcon(icon.NotFullData)
				}
				title := fmt.Sprintf("%s/%s%s", current, max, recovery)
				systray.SetTooltip(title)
				mResin.SetTitle(title)
				time.Sleep(10 * time.Second)
			}()
		}
	}()

	go func() {
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
