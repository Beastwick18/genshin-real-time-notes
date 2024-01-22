package main

import (
	"fmt"
	"resin/pkg/config"
	"resin/pkg/hoyo"
	"resin/pkg/hoyo/genshin"
	"resin/pkg/logging"
)

func main() {
	logging.SetFile("genshinDaily.log")

	const configFile string = "config.json"
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		logging.Fail("Failed loading config %s", configFile)
		return
	}
	json, err := hoyo.GetDailyData[genshin.GenshinDailyResponse](genshin.DailyURL, cfg.Ltoken, cfg.Ltuid, genshin.ActID)
	if err != nil {
		logging.Fail("Failed getting daily repsonse\n%s", err)
		return
	}
	fmt.Println(json.Message)
}
