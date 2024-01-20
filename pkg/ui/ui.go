package ui

import (
	"fmt"
	"image/color"
	"resin/pkg/config"
	"resin/pkg/db"
	"resin/pkg/logging"
	"time"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Beastwick18/go-webview2"
	"github.com/energye/systray"
	"github.com/skratchdot/open-golang/open"
)

type CommonMenu struct {
	Logs    *systray.MenuItem
	Refresh *systray.MenuItem
	Config  *systray.MenuItem
	Import  *systray.MenuItem
	Firefox *systray.MenuItem
	Quit    *systray.MenuItem
}

func CreateMenuItem(title string, icon []byte) *systray.MenuItem {
	item := systray.AddMenuItem(title, "")
	item.SetIcon(icon)
	return item
}

func createPopup(popup func(webview2.WebView, *config.Config), cfg *config.Config) {
	w := webview2.NewWithUserAgent(webview2.WebViewOptions{
		Debug:     true,
		AutoFocus: true,
		WindowOptions: webview2.WindowOptions{
			Title:  "Webview",
			PosX:   -404,
			PosY:   -745,
			Width:  384,
			Height: 654,
			IconId: 2, // icon resource id
			Center: false,
		},
	}, "Mozilla/5.0 (Linux; Android 11; SAMSUNG SM-G973U) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/14.2 Chrome/87.0.4280.141 Mobile Safari/537.36")
	if w == nil {
		logging.Fail("Failed to load webview")
		return
	}
	logging.Info("Opening webview")
	w.SetSize(384, 654, webview2.HintNone)
	popup(w, cfg)

	w.Run()
	w.Destroy()
}

func refreshLoop[T any](cfg *config.Config, menu *T, refresh func(*config.Config, *T)) {
	for {
		refresh(cfg, menu)
		time.Sleep(time.Duration(cfg.RefreshInterval) * time.Second)
	}
}

func watchEvents[T any](cm *CommonMenu, cfg *config.Config, menu *T, logFile string, configFile string, popup func(webview2.WebView, *config.Config), refresh func(*config.Config, *T)) {
	cm.Quit.Click(func() {
		systray.Quit()
	})
	cm.Refresh.Click(func() {
		logging.Info("User clicked refresh")
		refresh(cfg, menu)
	})
	cm.Logs.Click(func() {
		logging.Info(fmt.Sprintf("Opening \"%s\"", logFile))
		open.Start(logFile)
	})
	// cm.Config.Click(func() {
	// 	go func() {
	// 		cm.Config.Disable()
	// 		w := app.NewWindow()
	// 		// w.Option(app.Decorated(false))
	// 		w.Option(app.Size(500, 250))
	// 		err := run(w, cfg)
	// 		if err != nil {
	// 			logging.Fail(err.Error())
	// 		}
	// 		cm.Config.Enable()
	//
	// 	}()
	// })
	cm.Firefox.Click(func() {
		cookies := db.ReadFirefoxCookies()
		if cookies != nil {
			cfg.Ltoken = cookies.Ltoken_v2
			cfg.Ltuid = cookies.Ltuid_v2
			logging.Info("Got ltoken and ltuid from firefox")
			config.WriteConfig(cfg, configFile)
			refresh(cfg, menu)
		}
	})
	systray.SetOnClick(func(menu systray.IMenu) {
		createPopup(popup, cfg)
	})
}

func run(w *app.Window, cfg *config.Config) error {
	th := material.NewTheme()

	var ops op.Ops

	var ltoken, ltuid, genshinUID, hsrUID, refreshInterval widget.Editor

	ltoken.SetText(cfg.Ltoken)
	ltuid.SetText(cfg.Ltuid)
	genshinUID.SetText(cfg.GenshinUID)
	hsrUID.SetText(cfg.HsrUID)
	refreshInterval.SetText(fmt.Sprintf("%d", cfg.RefreshInterval))
	for {
		switch e := w.NextEvent().(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceEnd,
			}.Layout(gtx,
				generateInput(th, &ltoken, "ltoken", true),
				generateInput(th, &ltuid, "ltuid", true),
				generateInput(th, &genshinUID, "Genshin UID", false),
				generateInput(th, &hsrUID, "Honkai: Star Rail UID", false),
				generateInput(th, &refreshInterval, "Refresh Interval", false),
			)

			e.Frame(gtx.Ops)
		}
	}
}

func InitApp[T any](title string, tooltip string, icon []byte, logFile string, configFile string, menu *T, popup func(webview2.WebView, *config.Config), refresh func(*config.Config, *T)) {
	systray.SetOnClick(func(menu systray.IMenu) {
		menu.ShowMenu()
	})
	logging.SetFile(logFile)
	logging.Info("Application start")

	systray.SetIcon(icon)
	systray.SetTitle(title)
	systray.SetTooltip(tooltip)

	systray.AddSeparator()

	cm := &CommonMenu{}
	cm.Logs = systray.AddMenuItem("Logs", "Show logs")
	cm.Refresh = systray.AddMenuItem("Refresh", "Refresh data")
	// cm.Config = systray.AddMenuItem("Config", "Change the config")
	cm.Import = systray.AddMenuItem("Import Cookies", "Import cookies from browser")
	cm.Firefox = cm.Import.AddSubMenuItem("Firefox", "Load cookies from Firefox")
	cm.Quit = systray.AddMenuItem("Quit", "Exit the application")

	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		logging.Fail("Failed loading config file. Make sure it is present in the same directory you are running the program from.\n%s", err)
		systray.SetTooltip("Error loading config!")
	} else {
		go refreshLoop(cfg, menu, refresh)
	}

	go watchEvents(cm, cfg, menu, logFile, configFile, popup, refresh)
}

func generateInput(th *material.Theme, w *widget.Editor, hint string, mask bool) layout.FlexChild {
	return layout.Rigid(
		func(gtx layout.Context) layout.Dimensions {
			ed := material.Editor(th, w, hint)

			w.SingleLine = true
			w.Alignment = text.Start
			if mask {
				w.Mask = '●'
			}

			margins := layout.Inset{
				Top:   unit.Dp(10),
				Left:  unit.Dp(10),
				Right: unit.Dp(10),
			}
			padding := layout.Inset{
				Top:    unit.Dp(2),
				Bottom: unit.Dp(2),
				Left:   unit.Dp(5),
				Right:  unit.Dp(5),
			}
			// ... and borders ...
			border := widget.Border{
				Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
				CornerRadius: unit.Dp(3),
				Width:        unit.Dp(2),
			}

			// ... before laying it out, one inside the other
			return margins.Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					return border.Layout(gtx,
						func(gtx layout.Context) layout.Dimensions {
							return padding.Layout(gtx, ed.Layout)
						},
					)
				},
			)
		},
	)
}
