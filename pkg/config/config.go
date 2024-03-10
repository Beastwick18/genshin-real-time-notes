package config

import (
	"encoding/json"
	"io"
	"os"
	"resin/pkg/logging"
)

var VERSION string = "v0.0.4"

type Config struct {
	RefreshInterval int    `json:"refresh_interval"`
	UID             string `json:"uid"`
	Ltoken          string `json:"ltoken_v2"`
	Ltuid           string `json:"ltuid_v2"`
}

func LoadJSON[T any](reader io.Reader) (*T, error) {
	var cfg T
	bytesValue, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(bytesValue, &cfg)
	return &cfg, nil
}

func WriteConfig(cfg *Config, configPath string) error {
	var bt []byte
	var err error
	if bt, err = json.MarshalIndent(cfg, "", "    "); err != nil {
		return err
	}
	if err = os.WriteFile(configPath, bt, 0755); err != nil {
		return err
	}
	return nil
}

func LoadConfig(configPath string) (*Config, error) {
	if _, err := os.Stat(configPath); err != nil {
		logging.Fail("Unable to read config %s", configPath)
		return nil, err
	}
	var cfg Config
	jsonFile, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	bytesValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(bytesValue, &cfg)

	// Ensure at least one second of wait time before refresh
	cfg.RefreshInterval = max(1, cfg.RefreshInterval)

	return &cfg, nil
}
