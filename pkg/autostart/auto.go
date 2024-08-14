package autostart

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jxeng/shortcut"
)

type App struct {
	FileName         string
	Name             string
	Exec             []string
	WorkingDirectory string
}

func (a *App) IsEnabled() bool {
	startupPath := filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming", "Microsoft", "Windows", "Start Menu", "Programs", "Startup")
	savePath := filepath.Join(startupPath, a.FileName)
	_, err := os.Stat(savePath)
	return err == nil
}

func (a *App) Enable() error {
	startupPath := filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming", "Microsoft", "Windows", "Start Menu", "Programs", "Startup")
	savePath := filepath.Join(startupPath, a.FileName)
	exec := a.Exec[0]
	args := strings.Join(a.Exec[1:], " ")

	if err := os.MkdirAll(startupPath, 0777); err != nil {
		return err
	}
	sc := shortcut.Shortcut{
		ShortcutPath:     savePath,
		Target:           exec,
		IconLocation:     "",
		Arguments:        args,
		Description:      "",
		Hotkey:           "",
		WindowStyle:      "1",
		WorkingDirectory: a.WorkingDirectory,
	}
	return shortcut.Create(sc)
}

func (a *App) Disable() error {
	startupPath := filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming", "Microsoft", "Windows", "Start Menu", "Programs", "Startup")
	savePath := filepath.Join(startupPath, a.FileName)
	return os.Remove(savePath)
}
