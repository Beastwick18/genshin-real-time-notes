package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"resin/config"
	"resin/genshin"
	"resin/icon"
	"strconv"
	"time"

	"github.com/getlantern/systray"
)

func main() {
	systray.Run(onReady, onExit)
}

func refreshData(cfg config.Config, mResin *systray.MenuItem, mCommision *systray.MenuItem) {
	url := fmt.Sprintf("https://bbs-api-os.hoyolab.com/game_record/genshin/api/dailyNote?server=%s&role_id=%s", cfg.Server, cfg.Genshin_uuid)

	// Create a HTTP request
	rOrig := genshin.GenerateRequest(url, cfg.Ltoken, cfg.Ltuid)

	client := &http.Client{}
	for {
		func() {
			r := rOrig
			r.Header.Add("DS", genshin.GenerateDS())
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

			gr := genshin.LoadJSON(reader)

			current := gr.Data.CurrentResin
			max := gr.Data.MaxResin

			full := current == max
			seconds, _ := strconv.Atoi(gr.Data.ResinRecoveryTime)
			recovery := ""
			if seconds != 0 {
				hours := seconds / 3600
				minutes := (seconds / 60) - hours*60
				recovery = fmt.Sprintf(" [%dh %dm]", hours, minutes)
			}

			if full {
				systray.SetIcon(icon.FullData)
			} else {
				systray.SetIcon(icon.NotFullData)
			}
			title := fmt.Sprintf("%d/%d%s", current, max, recovery)
			systray.SetTooltip(title)
			mResin.SetTitle(title)
			mCommision.SetTitle(fmt.Sprintf("Commissions: %d/%d", gr.Data.FinishedTaskNum, gr.Data.TotalTaskNum))
			time.Sleep(10 * time.Second)
		}()
	}
}

func onReady() {
	systray.SetIcon(icon.NotFullData)
	systray.SetTitle("Genshin Real-Time Notes")
	systray.SetTooltip("?/?")
	mResin := systray.AddMenuItem("Resin: ?/?", "Resin")
	mResin.SetIcon(icon.NotFullData)

	mCommision := systray.AddMenuItem("Commissions: ?/?", "Commisions")

	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	cfg := config.LoadConfig("./config.json")

	go refreshData(cfg, mResin, mCommision)

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
