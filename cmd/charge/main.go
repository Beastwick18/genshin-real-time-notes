package main

import (
	"fmt"
	"resin/cmd"
	"resin/embedded"
	"resin/pkg/config"
	"resin/pkg/hoyo"
	"resin/pkg/hoyo/zzz"
	"resin/pkg/logging"
	"resin/pkg/ui"

	"github.com/energye/systray"
)

var logFile string = ".\\charge.log"
var configFile string = ".\\zzz_cookie.json"

type ZzzAssets struct {
	ChargeFull     []byte `asset:"zzz/charge_full.ico"`
	ChargeNotFull  []byte `asset:"zzz/charge_not_full.ico"`
	Engagement     []byte `asset:"zzz/engagement.ico"`
	EngagementDone []byte `asset:"zzz/engagement_done.ico"`
	CheckIn        []byte `asset:"zzz/checkin.ico"`
	Ticket         []byte `asset:"zzz/ticket.ico"`
	Tape           []byte `asset:"zzz/tape.ico"`
}

var assets ZzzAssets

type Menu struct {
	Charge      *systray.MenuItem
	Engagement  *systray.MenuItem
	ScratchCard *systray.MenuItem
	VideoStore  *systray.MenuItem
	CheckIn     *systray.MenuItem
}

var SaleStates = map[string]string{
	"SaleStateDoing": "Open",
	"SaleStateNo":    "Closed",
	"SaleStateDone":  "Done",
}

var CardSigns = map[string]string{
	"CardSignNo":   "Incomplete",
	"CardSignDone": "Done",
}

func refreshData(cfg *config.Config, m *Menu) {
	// First digit is always 1, differentiate by second digit
	server, ok := zzz.Servers[cfg.UID[1]]
	if !ok {
		logging.Fail(`Failed getting UID region (UID="%s")`, cfg.UID)
		systray.SetTooltip("Failed getting UID region!")
		return
	}
	zr, err := hoyo.GetData[zzz.ZzzResponse](zzz.BaseURL, server, cfg.UID, cfg.Ltoken, cfg.Ltuid)
	if err != nil {
		logging.Fail("Failed getting data from %s: Check your UID, ltoken, and ltuid\n%s", zzz.BaseURL, err)
		systray.SetTooltip("Failed getting data!")
		return
	}
	if zr.Retcode != 0 {
		logging.Fail("Server responded with (%d): %s\nCheck your UID, ltoken, and ltuid", zr.Retcode, zr.Message)
		systray.SetTooltip("Bad response from server!")
		systray.SetIcon(assets.ChargeFull)
		return
	}

	current := zr.Data.Energy.Progress.Current
	max := zr.Data.Energy.Progress.Max

	seconds := zr.Data.Energy.Restore
	recovery := ""
	if seconds != 0 {
		hours, minutes := hoyo.GetTime(seconds)
		recovery = fmt.Sprintf(" [%dh %dm]", hours, minutes)
	}

	charge := assets.ChargeNotFull
	if current == max {
		charge = assets.ChargeFull
	}
	systray.SetIcon(charge)
	m.Charge.SetIcon(charge)

	daily_current := zr.Data.Vitality.Current
	daily_max := zr.Data.Vitality.Max
	m.Engagement.SetTitle(fmt.Sprintf("Engagement: %d/%d", daily_current, daily_max))
	if daily_current == daily_max {
		m.Engagement.Disable()
		m.Engagement.SetIcon(assets.EngagementDone)
	} else {
		m.Engagement.Enable()
		m.Engagement.SetIcon(assets.Engagement)
	}

	saleState, ok := SaleStates[zr.Data.VhsSale.SaleState]
	if !ok {
		logging.Fail(fmt.Sprintf(`Failed to read video store state (state="%s")`, zr.Data.VhsSale.SaleState))
		m.VideoStore.SetTitle("Video Store: ERROR")
	} else {
		m.VideoStore.SetTitle(fmt.Sprintf("Video Store: %s", saleState))
	}

	scratchState, ok := CardSigns[zr.Data.CardSign]
	if !ok {
		logging.Fail(fmt.Sprintf(`Failed to read scratch card state (state="%s")`, zr.Data.CardSign))
		m.ScratchCard.SetTitle("Scratch Card: ERROR")
	} else {
		m.ScratchCard.SetTitle(fmt.Sprintf("Scratch Card: %s", scratchState))
	}

	title := fmt.Sprintf("%d/%d%s", current, max, recovery)
	systray.SetTooltip(title)
	m.Charge.SetTitle(title)
}

func checkIn(cfg *config.Config) {
	json, err := hoyo.GetDailyData[zzz.ZzzDailyResponse](zzz.DailyURL, cfg.Ltoken, cfg.Ltuid, zzz.ActID)
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
	m.Charge = ui.CreateMenuItem("Charge: ?/?", assets.ChargeNotFull)
	m.Engagement = ui.CreateMenuItem("Engagement: ?/?", assets.Engagement)
	m.ScratchCard = ui.CreateMenuItem("Scratch Card: ???", assets.Ticket)
	m.VideoStore = ui.CreateMenuItem("Video Store: ???", assets.Tape)
	m.CheckIn = ui.CreateMenuItem("Check In", assets.CheckIn)

	cfg := ui.InitApp("Zenless Zone Zero Real-Time Notes", "?/?", assets.ChargeNotFull, logFile, configFile, &m, "zzz", refreshData)
	watchEvents(cfg, &m)
}

func main() {
	cmd.ReadArgs(configFile, ".\\daily_zzz.log", checkIn)
	defer logging.CapturePanic()
	systray.Run(onReady, cmd.OnExit)
}
