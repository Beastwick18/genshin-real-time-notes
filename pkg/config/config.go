package config

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	Refresh_interval int    `json:"refresh_interval"`
	GenshinUID       string `json:"genshin_uid"`
	GenshinServer    string `json:"genshin_server"`
	HsrUID           string `json:"hsr_uid"`
	HsrServer        string `json:"hsr_server"`
	Ltoken           string `json:"ltoken"`
	Ltuid            string `json:"ltuid"`
}

var DefaultConfig Config = Config{
	Refresh_interval: 60,
	GenshinUID:       "genshin uid",
	GenshinServer:    "os_usa",
	HsrUID:           "hsr uid",
	HsrServer:        "prod_official_usa",
	Ltoken:           "token",
	Ltuid:            "mihoyo uid",
}

func LoadJSON[T any](reader io.Reader) T {
	var cfg T
	bytesValue, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(bytesValue, &cfg)
	return cfg
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
	json.Unmarshal(bytesValue, &cfg)

	// Ensure at least one second of wait time before refresh
	cfg.Refresh_interval = max(1, cfg.Refresh_interval)

	return &cfg
}
