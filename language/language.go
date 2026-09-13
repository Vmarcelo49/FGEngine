package language

import (
	"bytes"
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Language struct {
	Lang     Lang              `toml:"lang"`
	GameText map[string]string `toml:"game_text,omitempty"`
}

type Lang string

const (
	English    Lang = "EN"
	Portuguese Lang = "BR"
	Spanish    Lang = "SPA"

	defaultPath string = "./assets/text"
)

func LoadLang(configStr Lang) (*Language, error) {
	return ImportTOML(defaultPath + "/" + string(configStr) + ".toml")
}

func ImportTOML(filename string) (*Language, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var lang Language
	decoder := toml.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&lang); err != nil {
		return nil, fmt.Errorf("failed to decode language data: %w", err)
	}
	return &lang, nil
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
