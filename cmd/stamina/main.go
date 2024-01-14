package main

import (
	"fmt"
	"resin/pkg/config"
	"resin/pkg/hoyo"
	"resin/pkg/hoyo/hsr"
	"resin/pkg/icon"
	"resin/pkg/logging"
	"time"

	"github.com/getlantern/systray"
)

func refreshData(cfg *config.Config, mStamina *systray.MenuItem, mTraining *systray.MenuItem) {
	hr, err := hoyo.GetData[hsr.HsrResponse](hsr.BaseURL, cfg.HsrServer, cfg.HsrUID, cfg.Ltoken, cfg.Ltuid)
	if err != nil {
		logging.Fail("Failed getting data from %s: Check your UID, ltoken, and ltuid\n%s", hsr.BaseURL, err)
		systray.SetTooltip("Failed getting data!")
		return
	}

	current := hr.Data.CurrentStamina
	max := hr.Data.MaxStamina

	seconds := hr.Data.StaminaRecoveryTime
	recovery := ""
	if seconds != 0 {
		hours, minutes := hoyo.GetTime(seconds)
		recovery = fmt.Sprintf(" [%dh %dm]", hours, minutes)
	}

	if current == max {
		systray.SetIcon(icon.HsrFullData)
	} else {
		systray.SetIcon(icon.HsrNotFullData)
	}
	title := fmt.Sprintf("%d/%d%s", current, max, recovery)
	systray.SetTooltip(title)
	mStamina.SetTitle(title)
	mTraining.SetTitle(fmt.Sprintf("Training: %d/%d", hr.Data.CurrentTrainScore, hr.Data.MaxTrainScore))
}

func refreshDataLoop(cfg *config.Config, mStamina *systray.MenuItem, mTraining *systray.MenuItem) {
	for {
		refreshData(cfg, mStamina, mTraining)
		time.Sleep(time.Duration(cfg.Refresh_interval) * time.Second)
	}
}

func watchEvents(cfg *config.Config, mRefresh *systray.MenuItem, mQuit *systray.MenuItem, mStamina *systray.MenuItem, mTraining *systray.MenuItem) {
	for {
		select {
		case <-mQuit.ClickedCh:
			systray.Quit()
			return
		case <-mRefresh.ClickedCh:
			logging.Info("User clicked refresh")
			refreshData(cfg, mStamina, mTraining)
			break
		}
	}
}

func onReady() {
	logging.SetFile("./stamina.log")
	logging.Info("Application start")

	systray.SetIcon(icon.HsrNotFullData)
	systray.SetTitle("Honkai Star Rail Real-Time Notes")
	systray.SetTooltip("?/?")

	mStamina := systray.AddMenuItem("Stamina: ?/?", "Stamina")
	mStamina.SetIcon(icon.HsrNotFullData)

	mTraining := systray.AddMenuItem("Training: ?/?", "Training")

	systray.AddSeparator()

	mRefresh := systray.AddMenuItem("Refresh", "Refresh data")
	mQuit := systray.AddMenuItem("Quit", "Exit the application")

	cfg, err := config.LoadConfig("./config.json")
	if err != nil {
		logging.Fail("Failed loading config file. Make sure it is present in the same directory you are running the program from.\n%s", err)
		systray.SetTooltip("Error loading config!")
	} else {
		go refreshDataLoop(cfg, mStamina, mTraining)
	}

	go watchEvents(cfg, mRefresh, mQuit, mStamina, mTraining)
}

func onExit() {
	logging.Info("Exiting the application")
	logging.Close()
}

func main() {
	systray.Run(onReady, onExit)
}
