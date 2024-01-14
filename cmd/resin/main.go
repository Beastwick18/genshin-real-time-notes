package main

import (
	"fmt"
	"resin/pkg/config"
	"resin/pkg/hoyo"
	"resin/pkg/hoyo/genshin"
	"resin/pkg/icon"
	"resin/pkg/logging"
	"strconv"
	"time"

	"github.com/getlantern/systray"
)

func refreshData(cfg *config.Config, mResin *systray.MenuItem, mCommission *systray.MenuItem) {
	gr, err := hoyo.GetData[genshin.GenshinResponse](genshin.BaseURL, cfg.GenshinServer, cfg.GenshinUID, cfg.Ltoken, cfg.Ltuid)
	if err != nil {
		logging.Fail("Failed getting data from %s: Check your UID, ltoken, and ltuid\n%s", genshin.BaseURL, err)
		systray.SetTooltip("Failed getting data!")
		return
	}

	current := gr.Data.CurrentResin
	max := gr.Data.MaxResin

	seconds, err := strconv.Atoi(gr.Data.ResinRecoveryTime)

	var recovery string
	if err != nil {
		logging.Warn("Failed parsing recovery time (got \"%s\")\n%s", gr.Data.ResinRecoveryTime, err)
		recovery = " [?]"
	} else {
		hours, minutes := hoyo.GetTime(seconds)
		recovery = fmt.Sprintf(" [%dh %dm]", hours, minutes)
	}

	if current == max {
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
			return
		case <-mRefresh.ClickedCh:
			logging.Info("User clicked refresh")
			refreshData(cfg, mResin, mCommission)
			break
		}
	}
}

func onReady() {
	logging.SetFile("./resin.log")
	logging.Info("Application start")

	systray.SetIcon(icon.NotFullData)
	systray.SetTitle("Genshin Real-Time Notes")
	systray.SetTooltip("?/?")

	mResin := systray.AddMenuItem("Resin: ?/?", "Resin")
	mResin.SetIcon(icon.NotFullData)

	mCommission := systray.AddMenuItem("Commissions: ?/?", "Commissions")

	systray.AddSeparator()

	mRefresh := systray.AddMenuItem("Refresh", "Refresh data")
	mQuit := systray.AddMenuItem("Quit", "Exit the application")

	cfg, err := config.LoadConfig("./config.json")
	if err != nil {
		logging.Fail("Failed loading config file. Make sure it is present in the same directory you are running the program from.\n%s", err)
		systray.SetTooltip("Error loading config!")
	} else {
		go refreshDataLoop(cfg, mResin, mCommission)
	}

	go watchEvents(cfg, mRefresh, mQuit, mResin, mCommission)
}

func onExit() {
	logging.Info("Exiting the application")
	logging.Close()
}

func main() {
	systray.Run(onReady, onExit)
}
