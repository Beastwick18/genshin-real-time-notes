package ui

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"resin/pkg/config"
	"resin/pkg/logging"
	"time"

	"github.com/energye/systray"
	"github.com/skratchdot/open-golang/open"
)

type CommonMenu struct {
	Refresh  *systray.MenuItem
	Quit     *systray.MenuItem
	Advanced *systray.MenuItem
	Logs     *systray.MenuItem
	Login    *systray.MenuItem
}

func CreateMenuItem(title string, icon []byte) *systray.MenuItem {
	item := systray.AddMenuItem(title, "")
	item.SetIcon(icon)
	return item
}

func refreshLoop[T any](cfg *config.Config, menu *T, refresh func(*config.Config, *T)) {
	for {
		refresh(cfg, menu)
		logging.Info("Refreshed")
		time.Sleep(time.Duration(cfg.RefreshInterval) * time.Second)
	}
}

func watchEvents[T any](cm *CommonMenu, cfg *config.Config, menu *T, logFile string, configFile string, app string, refresh func(*config.Config, *T)) {
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
	cm.Login.Click(func() {
		var err error
		cfg, err = login(app, configFile, cfg, menu, refresh)
		if err != nil {
			logging.Fail("Failed to login:\n%s", err)
			return
		}
	})
}

func login[T any](app string, configFile string, cfg *config.Config, menu *T, refresh func(*config.Config, *T)) (*config.Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		logging.Fail("Failed to get working directory")
		return nil, err
	}
	exe := path.Join(wd, "login", "WinFormsApp1.exe")
	cmd := exec.Command(exe, app)
	cmd.Dir = "."
	// Block until finished
	_, err = cmd.CombinedOutput()
	if err != nil {
		logging.Fail("Failed to show login window:\n%s", err)
		return nil, err
	}
	logging.Info("Done")

	cookies, err := config.LoadConfig(configFile)
	if err != nil {
		logging.Fail("Failed to get webview cookies")
		return nil, err
	}
	logging.Info("Got ltoken and ltuid from webview")
	cfg = cookies
	refresh(cookies, menu)
	return cfg, nil
}

func InitApp[T any](title string, tooltip string, icon []byte, logFile string, configFile string, menu *T, app string, refresh func(*config.Config, *T)) *config.Config {
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

	cm.Advanced = systray.AddMenuItem("Advanced", "Advanced options")
	cm.Logs = cm.Advanced.AddSubMenuItem("Logs", "Show logs")
	cm.Login = cm.Advanced.AddSubMenuItem("Login", "Login To Hoyolab")

	cm.Refresh = systray.AddMenuItem("Refresh", "Refresh data")
	cm.Quit = systray.AddMenuItem("Quit", "Exit the application")

	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		cfg, err = login(app, configFile, cfg, menu, refresh)
		if err != nil {
			logging.Fail("Failed to login")
			return nil
		}
	}
	go refreshLoop(cfg, menu, refresh)

	watchEvents(cm, cfg, menu, logFile, configFile, app, refresh)
	return cfg
}
