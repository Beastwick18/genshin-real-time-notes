package ui

import (
	"fmt"
	"image/color"
	"resin/pkg/config"
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
	"github.com/energye/systray"
	"github.com/skratchdot/open-golang/open"
)

type CommonMenu struct {
	Logs    *systray.MenuItem
	Refresh *systray.MenuItem
	Config  *systray.MenuItem
	Quit    *systray.MenuItem
}

func CreateMenuItem(title string, icon []byte) *systray.MenuItem {
	item := systray.AddMenuItem(title, "")
	item.SetIcon(icon)
	return item
}

func refreshLoop[T any](cfg *config.Config, menu *T, refresh func(*config.Config, *T)) {
	for {
		refresh(cfg, menu)
		time.Sleep(time.Duration(cfg.RefreshInterval) * time.Second)
	}
}

func watchEvents[T any](cm *CommonMenu, cfg *config.Config, menu *T, logFile string, refresh func(*config.Config, *T)) {
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
	cm.Config.Click(func() {
		go func() {
			cm.Config.Disable()
			w := app.NewWindow()
			// w.Option(app.Decorated(false))
			w.Option(app.Size(500, 250))
			err := run(w, cfg)
			if err != nil {
				logging.Fail(err.Error())
			}
			cm.Config.Enable()
		}()
	})
}

func run(w *app.Window, cfg *config.Config) error {
	th := material.NewTheme()

	var ops op.Ops

	var ltoken widget.Editor
	var ltuid widget.Editor
	var genshinUID widget.Editor
	var genshinServer widget.Editor
	var hsrUID widget.Editor
	var hsrServer widget.Editor
	var refreshInterval widget.Editor
	ltoken.SetText(cfg.Ltoken)
	ltuid.SetText(cfg.Ltuid)
	genshinUID.SetText(cfg.GenshinUID)
	genshinServer.SetText(cfg.GenshinServer)
	hsrUID.SetText(cfg.HsrUID)
	hsrServer.SetText(cfg.HsrServer)
	refreshInterval.SetText(fmt.Sprintf("%d", cfg.RefreshInterval))
	for {
		switch e := w.NextEvent().(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			// This graphics context is used for managing the rendering state.
			gtx := layout.NewContext(&ops, e)

			layout.Flex{
				// Vertical alignment, from top to bottom
				Axis: layout.Vertical,
				// Empty space is left at the start, i.e. at the top
				Spacing: layout.SpaceEnd,
			}.Layout(gtx,
				generateInput(th, &ltoken, "ltoken", true),
				generateInput(th, &ltuid, "ltuid", true),
				generateInput(th, &genshinUID, "Genshin UID", false),
				generateInput(th, &genshinServer, "Genshin Server", false),
				generateInput(th, &hsrUID, "Honkai: Star Rail UID", false),
				generateInput(th, &hsrServer, "Honkai: Star Rail Server", false),
				generateInput(th, &refreshInterval, "Refresh Interval", false),
			)

			e.Frame(gtx.Ops)
		}
	}
}

func InitApp[T any](title string, tooltip string, icon []byte, logFile string, configFile string, menu *T, refresh func(*config.Config, *T)) {
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
	cm.Config = systray.AddMenuItem("Config", "Change the config")
	cm.Quit = systray.AddMenuItem("Quit", "Exit the application")

	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		logging.Fail("Failed loading config file. Make sure it is present in the same directory you are running the program from.\n%s", err)
		systray.SetTooltip("Error loading config!")
	} else {
		go refreshLoop(cfg, menu, refresh)
	}

	go watchEvents(cm, cfg, menu, logFile, refresh)
}

func generateInput(th *material.Theme, w *widget.Editor, hint string, mask bool) layout.FlexChild {
	return layout.Rigid(
		func(gtx layout.Context) layout.Dimensions {
			ed := material.Editor(th, w, hint)

			// Define characteristics of the input box
			w.SingleLine = true
			w.Alignment = text.Start
			if mask {
				w.Mask = '●'
			}

			// Define insets ...

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
