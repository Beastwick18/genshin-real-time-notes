package main

import (
	"fmt"
	"resin/cmd"
	"resin/embedded"
	"resin/pkg/config"
	"resin/pkg/hoyo"
	"resin/pkg/hoyo/genshin"
	"resin/pkg/logging"
	"resin/pkg/ui"
	"strconv"

	"github.com/Beastwick18/systray"
)

var logFile string = ".\\resin.log"
var configFile string = ".\\genshin_cookie.json"

type GenshinAssets struct {
	ResinFull    []byte `asset:"genshin/resin_full.ico"`
	ResinNotFull []byte `asset:"genshin/resin_not_full.ico"`
	Commission   []byte `asset:"genshin/commission.ico"`
	Expedition   []byte `asset:"genshin/expedition.ico"`
	Realm        []byte `asset:"genshin/realm.ico"`
	WeeklyBoss   []byte `asset:"genshin/weekly_boss.ico"`
	CheckIn      []byte `asset:"genshin/checkin.ico"`
}

var assets GenshinAssets

type Menu struct {
	Resin      *systray.MenuItem
	Commission *systray.MenuItem
	Expedition *systray.MenuItem
	Realm      *systray.MenuItem
	Domain     *systray.MenuItem
	CheckIn    *systray.MenuItem
}

func refreshData(cfg *config.Config, m *Menu) {
	if cfg == nil || m == nil {
		return
	}
	server := genshin.Servers[cfg.UID[0]]
	gr, err := hoyo.GetData[genshin.GenshinResponse](genshin.BaseURL, server, cfg.UID, cfg.Ltoken, cfg.Ltuid)
	if err != nil {
		logging.Fail("Failed getting data from %s: Check your UID, ltoken, and ltuid\n%s", genshin.BaseURL, err)
		systray.SetTooltip("Failed getting data!")
		return
	}
	if gr.Retcode != 0 {
		logging.Fail("Server responded with (%d): %s\nCheck your UID, ltoken, and ltuid", gr.Retcode, gr.Message)
		systray.SetTooltip("Bad response from server!")
		systray.SetIcon(assets.ResinFull)
		return
	}

	current := gr.Data.CurrentResin
	max := gr.Data.MaxResin

	seconds, err := strconv.Atoi(gr.Data.ResinRecoveryTime)

	var recovery string
	if err != nil {
		logging.Warn("Failed parsing recovery time (got \"%s\")\n%s", gr.Data.ResinRecoveryTime, err)
		recovery = " [?]"
	} else if seconds == 0 {
		recovery = ""
	} else {
		hours, minutes := hoyo.GetTime(seconds)
		recovery = fmt.Sprintf(" [%dh %dm]", hours, minutes)
	}

	if current == max {
		systray.SetIcon(assets.ResinFull)
	} else {
		systray.SetIcon(assets.ResinNotFull)
	}
	title := fmt.Sprintf("%d/%d%s", current, max, recovery)
	systray.SetTooltip(title)

	m.Resin.SetTitle(title)
	m.Commission.SetTitle(fmt.Sprintf("Commissions: %d/%d", gr.Data.FinishedTaskNum, gr.Data.TotalTaskNum))
	count := 0
	for _, exp := range gr.Data.Expeditions {
		if exp.Status == "Finished" {
			count++
		}
	}
	m.Expedition.SetTitle(fmt.Sprintf("Expeditions: %d/%d", count, gr.Data.MaxExpeditionNum))
	m.Realm.SetTitle(fmt.Sprintf("Realm: %d/%d", gr.Data.CurrentHomeCoin, gr.Data.MaxHomeCoin))
	m.Domain.SetTitle(fmt.Sprintf("Weekly Bosses: %d/%d", gr.Data.RemainResinDiscountNum, gr.Data.ResinDiscountNumLimit))
}

func checkIn(cfg *config.Config) {
	json, err := hoyo.GetDailyData[genshin.GenshinDailyResponse](genshin.DailyURL, cfg.Ltoken, cfg.Ltuid, genshin.ActID)
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

	err := embedded.ReadAssets(&assets)
	if err != nil {
		logging.Panic("Failed to read assets")
		return
	}

	var m Menu
	m.Resin = ui.CreateMenuItem("Resin: ?/?", assets.ResinNotFull)
	m.Commission = ui.CreateMenuItem("Commissions: ?/?", assets.Commission)
	m.Expedition = ui.CreateMenuItem("Expeditions: ?/?", assets.Expedition)
	m.Realm = ui.CreateMenuItem("Realm: ?/?", assets.Realm)
	m.Domain = ui.CreateMenuItem("Weekly Bosses: ?/?", assets.WeeklyBoss)
	m.CheckIn = ui.CreateMenuItem("Check In", assets.CheckIn)

	cfg := ui.InitApp("Genshin Real-Time Notes", "?/?", assets.ResinNotFull, logFile, configFile, &m, "genshin", refreshData)
	watchEvents(cfg, &m)
}

func main() {
	cmd.ReadArgs(configFile, ".\\daily_genshin.log", checkIn)
	defer logging.CapturePanic()
	systray.Run(onReady, cmd.OnExit)
}
