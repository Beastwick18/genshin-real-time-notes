package cmd

import (
	"fmt"
	"os"
	"resin/embedded"
	"resin/pkg/config"
	"resin/pkg/logging"
	"strings"
)

func OnExit() {
	logging.Info("Exiting the application")
	logging.Close()
	os.Exit(0)
}

func ReadArgs(configFile string, dailyFile string, checkIn func(*config.Config)) {
	args := os.Args[1:]
	no_extract := false
	for _, arg := range args {
		arg = strings.ToLower(arg)
		if arg == "--check-in" {
			logging.SetFile(dailyFile)
			cfg, err := config.LoadConfig(configFile)
			if err != nil {
				logging.Fail("Failed to load config")
				os.Exit(1)
			}
			checkIn(cfg)
			os.Exit(0)
		} else if arg == "-v" || arg == "--version" {
			fmt.Println(config.VERSION)
			os.Exit(0)
		} else if arg == "--no-extract" {
			no_extract = true
		}
	}
	if !no_extract {
		embedded.ExtractEmbeddedFiles()
	}
}
