package main

import (
	"fmt"
	"resin/pkg/config"
	"resin/pkg/hoyo"
	"resin/pkg/hoyo/genshin"
	"resin/pkg/icon"
	"strconv"
	"time"

	"github.com/getlantern/systray"
)

func refreshData(cfg *config.Config, mResin *systray.MenuItem, mCommission *systray.MenuItem) {
	gr := hoyo.GetData[genshin.GenshinResponse](genshin.BaseURL, cfg.GenshinServer, cfg.GenshinUID, cfg.Ltoken, cfg.Ltuid)

	current := gr.Data.CurrentResin
	max := gr.Data.MaxResin

	full := current == max
	seconds, _ := strconv.Atoi(gr.Data.ResinRecoveryTime)
	recovery := ""
	if seconds != 0 {
		hours, minutes := hoyo.GetTime(seconds)
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

func refreshDataLoop(cfg *config.Config, mResin *systray.MenuItem, mCommission *systray.MenuItem) {
	for {
		refreshData(cfg, mResin, mCommission)
		time.Sleep(time.Duration(cfg.Refresh_interval) * time.Second)
	}
}

func watchEvents(cfg *config.Config, mRefresh *systray.MenuItem, mQuit *systray.MenuItem, mResin *systray.MenuItem, mCommission *systray.MenuItem) {
	for {
		select {
		case <-mQuit.ClickedCh:
			systray.Quit()
			fmt.Println("quitting")
			return
		case <-mRefresh.ClickedCh:
			refreshData(cfg, mResin, mCommission)
			break
		}
	}
}

func onReady() {
	systray.SetIcon(icon.NotFullData)
	systray.SetTitle("Genshin Real-Time Notes")
	systray.SetTooltip("?/?")

	mResin := systray.AddMenuItem("Resin: ?/?", "Resin")
	mResin.SetIcon(icon.NotFullData)

	mCommission := systray.AddMenuItem("Commissions: ?/?", "Commissions")

	systray.AddSeparator()

	mRefresh := systray.AddMenuItem("Refresh", "Refresh data")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	cfg := config.LoadConfig("./config.json")

	go refreshDataLoop(cfg, mResin, mCommission)

	go watchEvents(cfg, mRefresh, mQuit, mResin, mCommission)
}

func onExit() {
	// clean up here
}

func main() {
	systray.Run(onReady, onExit)
}
