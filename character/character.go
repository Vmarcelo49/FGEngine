package character

import (
	"os"

	"gopkg.in/yaml.v2"

	"FGEngine/collision"
	"FGEngine/types"
)

type SpriteEx struct {
	ImagePath      string       `yaml:"imgPath"`
	Duration       uint         `yaml:"duration"`
	SourceSize     types.Rect   `yaml:"sourceSize"`
	Anchor         types.Rect   `yaml:"anchorPoint"`
	CollisionBoxes []types.Rect `yaml:"collisionBoxes"`
	HurtBoxes      []types.Rect `yaml:"hurtBoxes"`
	HitBoxes       []types.Rect `yaml:"hitBoxes"`
}

type Character struct {
	ID               int    `yaml:"id"`
	Name             string `yaml:"name"`
	Friction         int    `yaml:"friction"`
	JumpHeight       int    `yaml:"jumpHeight"`
	FilePath         string
	AnimationManager `yaml:"animationManager"`
	CharacterState   // info that matters when the game is running
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

// GetAllBoxes returns all boxes of the current character's sprite.
func (c *Character) GetAllBoxes() []collision.Box {
	var boxes []collision.Box
	for _, boxRect := range c.CurrentSprite.CollisionBoxes {
		boxes = append(boxes, collision.Box{Rect: boxRect, BoxType: collision.Collision})
	}
	for _, boxRect := range c.CurrentSprite.HitBoxes {
		boxes = append(boxes, collision.Box{Rect: boxRect, BoxType: collision.Hit})
	}
	for _, boxRect := range c.CurrentSprite.HurtBoxes {
		boxes = append(boxes, collision.Box{Rect: boxRect, BoxType: collision.Hurt})
	}
	return boxes
}
