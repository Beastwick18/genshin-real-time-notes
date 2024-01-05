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

var DefaultConfig Config = Config{
	Refresh_interval: 60,
	Genshin_uuid:     "genshin uuid",
	Server:           "os_usa",
	Ltoken:           "token",
	Ltuid:            "mihoyo uuid",
}

func WriteConfig(configPath string) {
	bt, _ := json.MarshalIndent(&DefaultConfig, "", "    ")
	os.WriteFile(configPath, bt, 0755)
}

func LoadConfig(configPath string) *Config {
	if _, err := os.Stat(configPath); err != nil {
		WriteConfig(configPath)
	}
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
	return &cfg
}
