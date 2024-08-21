package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	CommandJsonPath string `json:"commandJsonPath"`
	configPath string 
}

func GetConfig() Config {
	configPath := "./ypk.config.json"
	f, err := os.Open(configPath)
	if err != nil {
		return Config{
			CommandJsonPath: "./ypk.commands.json",
		}
	}
	defer f.Close()

	jsonParser := json.NewDecoder(f)
	var conf Config
	if err = jsonParser.Decode(&conf); err != nil {
		return Config{
			CommandJsonPath: "./ypk.commands.json",
		}
	}

	conf.configPath = configPath
	return conf
}