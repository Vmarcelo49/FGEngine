package animation

// Making a separate character package wouldnt be too useful, a character more or less just a place holder for animation data

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Character struct {
	ID         int                   `yaml:"id"`
	Name       string                `yaml:"name"`
	Friction   int                   `yaml:"friction"`
	JumpHeight int                   `yaml:"jumpHeight"`
	FilePath   string                // TODO, check if this is needed
	Animations map[string]*Animation `yaml:"animations"`
}

func LoadCharacter(path string) (*Character, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var character Character
	if err := yaml.Unmarshal(data, &character); err != nil {
		return nil, err
	}

	character.FilePath = path
	return &character, nil
}
