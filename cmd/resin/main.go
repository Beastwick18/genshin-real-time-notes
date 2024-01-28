package main

import (
	"fmt"
	"os"
	"resin/pkg/config"
	"resin/pkg/hoyo"
	"resin/pkg/hoyo/genshin"
	"resin/pkg/icon"
	"resin/pkg/logging"
	"resin/pkg/ui"
	"strconv"

	"github.com/energye/systray"
)

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
		systray.SetIcon(icon.FullData)
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
		systray.SetIcon(icon.FullData)
	} else {
		systray.SetIcon(icon.NotFullData)
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

func watchEvents(cfg *config.Config, m *Menu) {
	m.CheckIn.Click(func() {
		logging.Info("Clicked check in")
		json, err := hoyo.GetDailyData[genshin.GenshinDailyResponse](genshin.DailyURL, cfg.Ltoken, cfg.Ltuid, genshin.ActID)
		if err != nil {
			logging.Fail("Failed getting check in repsonse\n%s", err)
			return
		}
		logging.Info("%d: %s", json.Retcode, json.Message)
	})
}

func onReady() {
	defer logging.CapturePanic()
	m := &Menu{}
	m.Resin = ui.CreateMenuItem("Resin: ?/?", icon.NotFullData)
	m.Commission = ui.CreateMenuItem("Commissions: ?/?", icon.CommissionData)
	m.Expedition = ui.CreateMenuItem("Expeditions: ?/?", icon.ExpeditionData)
	m.Realm = ui.CreateMenuItem("Realm: ?/?", icon.RealmData)
	m.Domain = ui.CreateMenuItem("Weekly Bosses: ?/?", icon.WeeklyBossData)
	m.CheckIn = ui.CreateMenuItem("Check In", icon.GenshinCheckIn)

	cfg := ui.InitApp("Genshin Real-Time Notes", "?/?", icon.NotFullData, ".\\resin.log", ".\\genshin_cookie.json", m, "genshin", refreshData)
	watchEvents(cfg, m)
}

func onExit() {
	logging.Info("Exiting the application")
	logging.Close()
	os.Exit(0)
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			logging.Panic("%v", err)
		}
	}()
	systray.Run(onReady, onExit)
}
