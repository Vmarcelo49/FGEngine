package character

import (
	"fgengine/animation"
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
	FilePath   string                          `yaml:"filepath,omitempty"`
	Animations map[string]*animation.Animation `yaml:"animations"`

	// Ingame Related
	HP                  int                 `yaml:"hp,omitempty"`
	Position            types.Vector2       `yaml:"position,omitempty"`
	Velocity            types.Vector2       `yaml:"-"`
	IgnoreGravityFrames int                 `yaml:"-"`
	StateMachine        *state.StateMachine `yaml:"-"`

	// animation related
	ActiveAnimation *animation.Animation `yaml:"-"`
	ActiveSprite    *animation.Sprite    `yaml:"-"`
	FrameIndex      int                  `yaml:"-"`
	SpriteIndex     int                  `yaml:"-"`
	ShouldLoop      bool                 `yaml:"-"`
	AnimationQueue  []string             `yaml:"-"`
	AnimationIndex  int                  `yaml:"-"`
}

type CharacterID int

const (
	Helmet CharacterID = iota
)

// LoadCharacter loads a character by its ID.
// in the future, this and update should be the only two exported functions
func LoadCharacter(id CharacterID) (*Character, error) {
	chara := &Character{}
	var err error
	switch id {
	case Helmet:
		chara, err = loadCharacterByFile(defaultCharacterPath + "helmet.yaml")
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown character ID: %d", id)
	}
	chara.initialize()
	return chara, nil
}

func (c *Character) initialize() {
	c.setAnimation("idle")
	if c.ActiveAnimation == nil {
		panic("Character must have an 'idle' animation")
	}
	c.ActiveSprite = c.ActiveAnimation.Sprites[0]
	c.StateMachine = &state.StateMachine{}
}

func (c *Character) setAnimation(name string) {
	anim, ok := c.Animations[name]
	if !ok {
		return
	}
	c.ActiveAnimation = anim
	c.ShouldLoop = true
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

// Funcs for the Renderable interface:
func (c *Character) GetID() int {
	return c.ID
}

func (c *Character) GetPosition() types.Vector2 {
	return c.Position
}

func (c *Character) GetSprite() *animation.Sprite {
	return c.ActiveSprite
}
