package config

import (
	"bufio"
	"encoding/json"
	"os"
)

const CONFIG_FILE = "bridgeSettings.json"

type BridgeConfig struct {
    Ip     string `json:"ip"`
    ApiKey string `json:"apikey"`
}

func (c *BridgeConfig) Load() BridgeConfig {
	f, err := os.Open(CONFIG_FILE)
	if err != nil {
		return BridgeConfig{}
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	content := ""

	for scanner.Scan() {
		content += scanner.Text()
	}

	err = json.Unmarshal([]byte(content), c)
	if err != nil {
		return *c
	}

	return *c
}

func (c *BridgeConfig) Save() error {
	f, err := os.Create(CONFIG_FILE)
	if err != nil {
		return err
	}
	defer f.Close()

	configBytes, err := json.Marshal(*c)
	if err != nil {
		return err
	}

	f.Write(configBytes)

    return nil
}
