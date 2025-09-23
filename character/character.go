package character

import (
	"fgengine/animation"
	"fgengine/state"
	"fgengine/types"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	defaultCharacterPath = "./assets/characters/"
)

type Character struct {
	ID         int                             `yaml:"id"`
	Name       string                          `yaml:"name"`
	Friction   int                             `yaml:"friction"`
	JumpHeight int                             `yaml:"jumpHeight"`
	FilePath   string                          `yaml:"filepath,omitempty"`
	Animations map[string]*animation.Animation `yaml:"animations"`

	// Ingame Related
	HP                  int                 `yaml:"-"`
	Position            types.Vector2       `yaml:"-"`
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

	// Convert relative paths to absolute paths based on YAML file location
	for _, anim := range character.Animations {
		for _, sprite := range anim.Sprites {
			if sprite.ImagePath != "" {
				sprite.ImagePath = resolveRelativePath(sprite.ImagePath, filePath)
			}
		}
	}

	return &character, nil
}

// resolveRelativePath converts a relative path to an absolute path based on a reference path
func resolveRelativePath(relativePath, referencePath string) string {
	if filepath.IsAbs(relativePath) {
		return relativePath
	}
	referenceDir := filepath.Dir(referencePath)
	return filepath.Clean(filepath.Join(referenceDir, relativePath))
}

func (c *Character) GetPosition() types.Vector2 {
	return c.Position
}

func (c *Character) GetSprite() *animation.Sprite {
	return c.ActiveSprite
}
