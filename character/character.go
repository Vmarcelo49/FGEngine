package character

import (
	"fgengine/animation"
	"fgengine/collision"
	"fgengine/state"
	"fgengine/types"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	defaultCharacterPath = "./assets/characters/"
)

// well lets try hard to not cause a circular dependency with animation and player

type Character struct {
	ID         int                             `yaml:"id"`
	Name       string                          `yaml:"name"`
	Friction   int                             `yaml:"friction"`
	JumpHeight int                             `yaml:"jumpHeight"`
	FilePath   string                          // TODO, check if this is needed, probably only used in the editor?
	Animations map[string]*animation.Animation `yaml:"animations"`

	// Ingame Related
	HP                  int
	Position            types.Vector2
	Velocity            types.Vector2
	IgnoreGravityFrames int
	StateMachine        *state.StateMachine

	// animation related
	ActiveAnimation *animation.Animation
	activeSprite    *animation.Sprite
	FrameIndex      int
	SpriteIndex     int
	ShouldLoop      bool
	AnimationQueue  []string
	AnimationIndex  int
}

type CharacterID int

const (
	Helmet CharacterID = iota
)

func LoadCharacter(id CharacterID) (*Character, error) {
	switch id {
	case Helmet:
		return loadCharacterByFile(defaultCharacterPath + "helmet.yaml")
	default:
		return nil, fmt.Errorf("unknown character ID: %d", id)
	}
}

func loadCharacterByFile(filePath string) (*Character, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var character Character
	if err := yaml.Unmarshal(data, &character); err != nil {
		return nil, err
	}

	return &character, nil
}

func (c *Character) GetID() int {
	return c.ID
}

// Funcs for the Renderable interface:

func (c *Character) GetPosition() types.Vector2 {
	return c.Position
}

func (c *Character) GetAllBoxes() []collision.Box {
	return c.AnimationSystem.GetCurrentSprite().GetAllBoxes()
}

func (c *Character) GetSprite() *animation.Sprite {
	return c.AnimationSystem.GetCurrentSprite()
}
