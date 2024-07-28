package main

import (
	"fmt"
	"resin/cmd"
	"resin/embedded"
	"resin/pkg/config"
	"resin/pkg/hoyo"
	"resin/pkg/hoyo/hsr"
	"resin/pkg/logging"
	"resin/pkg/ui"

	"github.com/energye/systray"
)

var logFile string = ".\\stamina.log"
var configFile string = ".\\hsr_cookie.json"

type HsrAssets struct {
	StaminaFull    []byte `asset:"hsr/stamina_full.ico"`
	StaminaNotFull []byte `asset:"hsr/stamina_not_full.ico"`
	Training       []byte `asset:"hsr/training.ico"`
	Expedition     []byte `asset:"hsr/expedition.ico"`
	EchoOfWar      []byte `asset:"hsr/echo_of_war.ico"`
	CheckIn        []byte `asset:"hsr/checkin.ico"`
}

var assets HsrAssets

type Menu struct {
	Stamina    *systray.MenuItem
	Training   *systray.MenuItem
	Expedition *systray.MenuItem
	Reserve    *systray.MenuItem
	EchoOfWar  *systray.MenuItem
	CheckIn    *systray.MenuItem
}

func refreshData(cfg *config.Config, m *Menu) {
	server := hsr.Servers[cfg.UID[0]]
	hr, err := hoyo.GetData[hsr.HsrResponse](hsr.BaseURL, server, cfg.UID, cfg.Ltoken, cfg.Ltuid)
	if err != nil {
		logging.Fail("Failed getting data from %s: Check your UID, ltoken, and ltuid\n%s", hsr.BaseURL, err)
		systray.SetTooltip("Failed getting data!")
		return
	}
	if hr.Retcode != 0 {
		logging.Fail("Server responded with (%d): %s\nCheck your UID, ltoken, and ltuid", hr.Retcode, hr.Message)
		systray.SetTooltip("Bad response from server!")
		systray.SetIcon(assets.StaminaFull)
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
		systray.SetIcon(assets.StaminaFull)
	} else {
		systray.SetIcon(assets.StaminaNotFull)
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

func checkIn(cfg *config.Config) {
	json, err := hoyo.GetDailyData[hsr.HsrDailyResponse](hsr.DailyURL, cfg.Ltoken, cfg.Ltuid, hsr.ActID)
	if err != nil {
		logging.Fail("Failed getting check in repsonse\n%s", err)
		return
	}
	logging.Info("%d: %s", json.Retcode, json.Message)
}

func watchEvents(cfg *config.Config, m *Menu) {
	m.CheckIn.Click(func() {
		logging.Info("Clicked check in")
		checkIn(cfg)
	})
}

func onReady() {
	defer logging.CapturePanic()
	logging.SetFile(logFile)

	embedded.ReadAssets(&assets)

	m := &Menu{}
	m.Stamina = ui.CreateMenuItem("Stamina: ?/?", assets.StaminaNotFull)
	m.Training = ui.CreateMenuItem("Training: ?/?", assets.Training)
	m.Expedition = ui.CreateMenuItem("Expeditions: ?/?", assets.Expedition)
	m.Reserve = ui.CreateMenuItem("Reserve: ?/?", assets.StaminaFull)
	m.EchoOfWar = ui.CreateMenuItem("Echo of War: ?/?", assets.EchoOfWar)
	m.CheckIn = ui.CreateMenuItem("Check In", assets.CheckIn)

	cfg := ui.InitApp("Honkai Star Rail Real-Time Notes", "?/?", assets.StaminaNotFull, logFile, configFile, m, "hsr", refreshData)
	watchEvents(cfg, m)
}

func main() {
	cmd.ReadArgs(configFile, ".\\daily_hsr.log", checkIn)
	defer logging.CapturePanic()
	systray.Run(onReady, cmd.OnExit)
}
