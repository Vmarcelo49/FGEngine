package character

import (
	"os"

	"gopkg.in/yaml.v2"

	"FGEngine/types"
)

type Sprite struct {
	ImagePath      string       `yaml:"imgPath"`
	Duration       uint         `yaml:"duration"`
	SourceSize     types.Rect   `yaml:"sourceSize"`
	Anchor         types.Rect   `yaml:"anchorPoint"`
	CollisionBoxes []types.Rect `yaml:"collisionBoxes"`
	HurtBoxes      []types.Rect `yaml:"hurtBoxes"`
	HitBoxes       []types.Rect `yaml:"hitBoxes"`
}

type Character struct {
	ID         int    `yaml:"id"`
	Name       string `yaml:"name"`
	Friction   int    `yaml:"friction"`
	JumpHeight int    `yaml:"jumpHeight"`
	FilePath   string
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
