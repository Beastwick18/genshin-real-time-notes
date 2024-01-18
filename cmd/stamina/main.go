package main

import (
	"fmt"
	"os"
	"resin/pkg/config"
	"resin/pkg/hoyo"
	"resin/pkg/hoyo/hsr"
	"resin/pkg/icon"
	"resin/pkg/logging"
	"resin/pkg/ui"

	"gioui.org/app"
	"github.com/Beastwick18/go-webview2"
	"github.com/energye/systray"
)

type Menu struct {
	Stamina    *systray.MenuItem
	Training   *systray.MenuItem
	Expedition *systray.MenuItem
	Reserve    *systray.MenuItem
	EchoOfWar  *systray.MenuItem
}

func refreshData(cfg *config.Config, m *Menu) {
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

	m.Stamina.SetTitle(title)
	m.Training.SetTitle(fmt.Sprintf("Training: %d/%d", hr.Data.CurrentTrainScore, hr.Data.MaxTrainScore))
	count := 0
	for _, exp := range hr.Data.Expeditions {
		if exp.Status == "Finished" {
			count++
		}
	}
	m.Expedition.SetTitle(fmt.Sprintf("Expeditions: %d/%d", count, hr.Data.TotalExpeditionNum))
	m.Reserve.SetTitle(fmt.Sprintf("Reserve: %d/2400", hr.Data.CurrentReserveStamina))
	m.EchoOfWar.SetTitle(fmt.Sprintf("Echo of War: %d/%d", hr.Data.WeeklyCocoonCnt, hr.Data.WeeklyCocoonLimit))
}

func popup(w webview2.WebView, cfg *config.Config) {
	w.SetTitle("Honkai: Star Rail")
	w.Navigate("https://act.hoyolab.com/app/community-game-records-sea/rpg/m.html#/hsr")
}

func onReady() {
	m := &Menu{}
	m.Stamina = ui.CreateMenuItem("Stamina: ?/?", icon.HsrNotFullData)
	m.Training = ui.CreateMenuItem("Training: ?/?", icon.TrainingData)
	m.Expedition = ui.CreateMenuItem("Expeditions: ?/?", icon.HsrExpeditionData)
	m.Reserve = ui.CreateMenuItem("Expeditions: ?/?", icon.HsrFullData)
	m.EchoOfWar = ui.CreateMenuItem("Echo of War: ?/?", icon.EchoOfWarData)

	ui.InitApp("Honkai Star Rail Real-Time Notes", "?/?", icon.HsrNotFullData, ".\\stamina.log", ".\\config.json", m, popup, refreshData)
}

func onExit() {
	logging.Info("Exiting the application")
	logging.Close()
	os.Exit(0)
}

func main() {
	systray.Run(onReady, onExit)
	app.Main()
}
