package ui

import (
	"fmt"
	"resin/pkg/config"
	"resin/pkg/logging"
	"time"

	"github.com/energye/systray"
	"github.com/skratchdot/open-golang/open"
)

type CommonMenu struct {
	Logs    *systray.MenuItem
	Refresh *systray.MenuItem
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
		time.Sleep(time.Duration(cfg.Refresh_interval) * time.Second)
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
