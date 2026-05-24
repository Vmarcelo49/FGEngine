package config

// Only user-configurable settings should be here

import (
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"gopkg.in/yaml.v3"
)

const configFilePath = "config.yaml"

var ActiveConfig Config

type Config struct {
	WindowWidth        int     `yaml:"window_width"`
	WindowHeight       int     `yaml:"window_height"`
	ControllerDeadzone float64 `yaml:"deadzone"`
	Language           string  `yaml:"language"`
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
	ebiten.SetWindowSize(ActiveConfig.WindowWidth, ActiveConfig.WindowWidth)
	ebiten.SetWindowTitle("FG Engine")
}

// LoadConfigFile tries to load the config file found in the same path as the executable, if its not found one is created from the default config
func LoadConfigFile() Config {
	var config Config
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		config = loadDefaultConfig()
		yamlbytes, err := yaml.Marshal(config)
		if err != nil {
			log.Fatalf("error creating default config file: %s", err.Error())
		}
		file, err := os.Create(configFilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			err = file.Close()
		}()
		if err != nil {
			log.Fatal(err)
		}

		_, err = file.Write(yamlbytes)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		log.Fatalf("error loading config file: %s", err.Error())
	}
	return config
}
