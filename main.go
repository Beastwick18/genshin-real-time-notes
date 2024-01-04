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

func refreshData(r *http.Request, client *http.Client, mResin *systray.MenuItem, mCommission *systray.MenuItem) {
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
	mCommission.SetTitle(fmt.Sprintf("Commissions: %d/%d", gr.Data.FinishedTaskNum, gr.Data.TotalTaskNum))
}

func refreshDataLoop(cfg config.Config, mResin *systray.MenuItem, mCommission *systray.MenuItem) {
	// Create a HTTP request
	rOrig := genshin.GenerateRequest(cfg.Server, cfg.Genshin_uuid, cfg.Ltoken, cfg.Ltuid)

	client := &http.Client{}
	for {
		refreshData(rOrig, client, mResin, mCommission)
		time.Sleep(time.Duration(cfg.Refresh_interval) * time.Second)
	}
}

func onReady() {
	systray.SetIcon(icon.NotFullData)
	systray.SetTitle("Genshin Real-Time Notes")
	systray.SetTooltip("?/?")
	mResin := systray.AddMenuItem("Resin: ?/?", "Resin")
	mResin.SetIcon(icon.NotFullData)

	mCommission := systray.AddMenuItem("Commissions: ?/?", "Commisions")

	systray.AddSeparator()
	mRefresh := systray.AddMenuItem("Refresh", "Refresh data")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	cfg := config.LoadConfig("./config.json")

	go refreshDataLoop(cfg, mResin, mCommission)

	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
				fmt.Println("quitting")
				return
			case <-mRefresh.ClickedCh:
				rOrig := genshin.GenerateRequest(cfg.Server, cfg.Genshin_uuid, cfg.Ltoken, cfg.Ltuid)
				client := &http.Client{}
				refreshData(rOrig, client, mResin, mCommission)
				break
			}
		}
	}()
}

func onExit() {
	// clean up here
}
