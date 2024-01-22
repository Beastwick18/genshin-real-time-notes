package main

import (
	"fmt"
	"resin/pkg/config"
	"resin/pkg/hoyo"
	"resin/pkg/hoyo/hsr"
	"resin/pkg/logging"
)

func main() {
	logging.SetFile("hsrDaily.log")

	const configFile string = "config.json"
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		logging.Fail("Failed loading config %s", configFile)
		return
	}
	json, err := hoyo.GetDailyData[hsr.HsrDailyResponse](hsr.DailyURL, cfg.Ltoken, cfg.Ltuid, hsr.ActID)
	if err != nil {
		logging.Fail("Failed getting daily repsonse\n%s", err)
		return
	}
	fmt.Println(json.Message)
}
