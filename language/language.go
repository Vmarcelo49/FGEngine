package language

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Language struct {
	Lang     Lang              `yaml:"lang"`
	GameText map[string]string `yaml:"game_text"`
}

type Lang string

const (
	English    Lang = "EN"
	Portuguese Lang = "BR"
	Spanish    Lang = "SPA"

	defaultPath string = "./assets/text"
)

func LoadLang(configStr Lang) (*Language, error) {
	lang, err := ImportYAML(defaultPath + string(configStr))
	if err != nil {
		return nil, err
	}
	return lang, nil
}

func (lang *Language) exportYAML(filename string) error {
	data, err := yaml.Marshal(lang)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func ImportYAML(filename string) (*Language, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var lang Language
	err = yaml.Unmarshal(data, &lang)
	return &lang, err
}

func MakePTBR() Language {
	return Language{
		GameText: map[string]string{
			"match":  "Jogar",
			"config": "Configuração",
			"exit":   "Sair",
		},
		Lang: Portuguese,
	}
}
