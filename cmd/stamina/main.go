package main

import (
	"fmt"
	"resin/pkg/config"
	"resin/pkg/hoyo"
	"resin/pkg/hoyo/hsr"
	"resin/pkg/icon"
	"time"

	"github.com/getlantern/systray"
)

func refreshData(cfg *config.Config, mStamina *systray.MenuItem, mTraining *systray.MenuItem) {
	hr := hoyo.GetData[hsr.HsrResponse](hsr.BaseURL, cfg.HsrServer, cfg.HsrUID, cfg.Ltoken, cfg.Ltuid)

	current := hr.Data.CurrentStamina
	max := hr.Data.MaxStamina

	seconds := hr.Data.StaminaRecoveryTime
	recovery := ""
	if seconds != 0 {
		hours := seconds / 3600
		minutes := (seconds / 60) - hours*60
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
			fmt.Println("quitting")
			return
		case <-mRefresh.ClickedCh:
			refreshData(cfg, mStamina, mTraining)
			break
		}
	}
}

func onReady() {
	systray.SetIcon(icon.HsrNotFullData)
	systray.SetTitle("Honkai Star Rail Real-Time Notes")
	systray.SetTooltip("?/?")

	mStamina := systray.AddMenuItem("Stamina: ?/?", "Stamina")
	mStamina.SetIcon(icon.HsrNotFullData)

	mTraining := systray.AddMenuItem("Training: ?/?", "Training")

	systray.AddSeparator()

	mRefresh := systray.AddMenuItem("Refresh", "Refresh data")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	cfg := config.LoadConfig("./config.json")

	go refreshDataLoop(cfg, mStamina, mTraining)

	go watchEvents(cfg, mRefresh, mQuit, mStamina, mTraining)
}

func onExit() {
	// clean up here
}

func main() {
	systray.Run(onReady, onExit)
}
