package config

// Only user-configurable settings should be here

import (
	"bytes"
	"log"
	"os"

	"github.com/pelletier/go-toml/v2"
)

const configFilePath = "config.toml"

var ActiveConfig Config

type Config struct {
	WindowWidth        int     `toml:"window_width"`
	WindowHeight       int     `toml:"window_height"`
	ControllerDeadzone float64 `toml:"deadzone"`
	Language           string  `toml:"language"`
}

func loadDefaultConfig() Config {
	return Config{
		WindowWidth:        1600,
		WindowHeight:       900,
		ControllerDeadzone: 0.3,
		Language:           "EN",
	}
}

func InitGameConfig() {
	ActiveConfig = LoadConfigFile()
}

// LoadConfigFile tries to load the config file found in the same path as where the project was ran, if its not found one is created from the default config
func LoadConfigFile() Config {
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		config := loadDefaultConfig()
		tomlbytes, err := toml.Marshal(config)
		if err != nil {
			log.Fatalf("error creating default config file: %s", err.Error())
		}
		if err := os.WriteFile(configFilePath, tomlbytes, 0644); err != nil {
			log.Fatal(err)
		}
		return config
	}

	var config Config
	decoder := toml.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&config); err != nil {
		log.Fatalf("error loading config file: %s", err.Error())
	}
	return config
}
