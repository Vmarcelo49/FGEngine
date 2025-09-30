package character

import (
	"fgengine/animation"
	"fgengine/constants"
	"fgengine/graphics"
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
	StateMachine *state.StateMachine `yaml:"-"`

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
	return c.StateMachine.Position
}

func (c *Character) GetSprite() *animation.Sprite {
	return c.ActiveSprite
}

func (c *Character) GetRenderProperties() graphics.RenderProperties {
	// For now, return default properties. In the future, you could add
	// character-specific properties like scale for different sized characters,
	// layer for draw order, or color modulation for effects
	return graphics.DefaultRenderProperties()
}

func (c *Character) Update() {
	switch c.StateMachine.ActiveState {
	case state.StateIdle:
		c.setAnimation("idle")
	case state.StateWalk | state.StateForward:
		c.setAnimation("walk")
		c.StateMachine.Velocity.X = 2 // in the future, this should be based on something we can get from the editor
	case state.StateDash | state.StateForward:
		c.setAnimation("dash")
		c.StateMachine.Velocity.X = 5
	default:
		c.setAnimation("idle")
	}
	// Update position based on velocity
	c.StateMachine.Position.X += c.StateMachine.Velocity.X
	c.StateMachine.Position.Y += c.StateMachine.Velocity.Y

	// Apply friction
	if c.StateMachine.Velocity.X > 0 {
		c.StateMachine.Velocity.X -= float64(c.Friction)
		if c.StateMachine.Velocity.X < 0 {
			c.StateMachine.Velocity.X = 0
		}
	} else if c.StateMachine.Velocity.X < 0 {
		c.StateMachine.Velocity.X += float64(c.Friction)
		if c.StateMachine.Velocity.X > 0 {
			c.StateMachine.Velocity.X = 0
		}
	}

	// Apply gravity only if airborne and not ignoring gravity frames
	if c.StateMachine.HasState(state.StateAirborne) && c.StateMachine.IgnoreGravityFrames <= 0 {
		c.StateMachine.Velocity.Y += constants.Gravity // example gravity value
	} else {
		c.StateMachine.IgnoreGravityFrames--
	}

	if c.StateMachine.Position.Y >= constants.GroundLevelY {
		c.StateMachine.Position.Y = constants.GroundLevelY
		c.StateMachine.Velocity.Y = 0
		c.StateMachine.AddState(state.StateGrounded)
		c.StateMachine.RemoveState(state.StateAirborne)
	}

	if c.StateMachine.Position.X < constants.World.X {
		c.StateMachine.Position.X = constants.World.X
	}
	if c.StateMachine.Position.X+float64(c.ActiveSprite.Rect.W) > constants.World.Right() {
		c.StateMachine.Position.X = constants.World.Right() - float64(c.ActiveSprite.Rect.W)
	}
}
