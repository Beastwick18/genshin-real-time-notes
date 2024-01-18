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

	"gioui.org/app"
	"github.com/Beastwick18/go-webview2"
	"github.com/energye/systray"
)

type Menu struct {
	Resin      *systray.MenuItem
	Commission *systray.MenuItem
	Expedition *systray.MenuItem
	Realm      *systray.MenuItem
	Domain     *systray.MenuItem
}

func refreshData(cfg *config.Config, m *Menu) {
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

func popup(menu systray.IMenu, cfg *config.Config) {
	w := webview2.NewWithUserAgent(webview2.WebViewOptions{
		Debug:     true,
		AutoFocus: true,
		WindowOptions: webview2.WindowOptions{
			Title:  "Genshin",
			PosX:   -404,
			PosY:   -745,
			Width:  384,
			Height: 654,
			IconId: 2, // icon resource id
			Center: false,
		},
	}, "Mozilla/5.0 (Linux; Android 11; SAMSUNG SM-G973U) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/14.2 Chrome/87.0.4280.141 Mobile Safari/537.36")
	if w == nil {
		logging.Fail("Failed to load webview.")
		return
	}
	w.SetSize(384, 654, webview2.HintNone)
	w.Navigate(fmt.Sprintf("https://act.hoyolab.com/app/community-game-records-sea/m.html#/ys/realtime?role_id=%s&server=%s", cfg.GenshinUID, cfg.GenshinServer))

	w.Run()
	w.Destroy()
}

func onReady() {
	m := &Menu{}
	m.Resin = ui.CreateMenuItem("Resin: ?/?", icon.NotFullData)
	m.Commission = ui.CreateMenuItem("Commissions: ?/?", icon.CommissionData)
	m.Expedition = ui.CreateMenuItem("Expeditions: ?/?", icon.ExpeditionData)
	m.Realm = ui.CreateMenuItem("Realm: ?/?", icon.RealmData)
	m.Domain = ui.CreateMenuItem("Weekly Bosses: ?/?", icon.WeeklyBossData)

	ui.InitApp("Genshin Real-Time Notes", "?/?", icon.NotFullData, ".\\resin.log", ".\\config.json", m, popup, refreshData)
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
