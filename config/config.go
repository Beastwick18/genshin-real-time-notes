package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	Refresh_interval int    `json:"refresh_interval"`
	Genshin_uuid     string `json:"genshin_uuid"`
	Server           string `json:"server"`
	Ltoken           string `json:"ltoken"`
	Ltuid            string `json:"ltuid"`
}

func LoadConfig(configPath string) Config {
	var cfg Config
	jsonFile, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	bytesValue, err := io.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}
	fmt.Println(bytesValue)
	json.Unmarshal(bytesValue, &cfg)
	return cfg
}
