package character

import (
	"fgengine/animation"
	"fgengine/state"
	"fgengine/types"
	"fmt"
	"image"
	_ "image/png"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	defaultCharacterPath = "./assets/characters/"
)

type Character struct {
	ID         int                             `yaml:"id,omitempty"`
	Name       string                          `yaml:"name"`
	Friction   float64                         `yaml:"friction,omitempty"`
	JumpHeight float64                         `yaml:"jumpHeight,omitempty"`
	FilePath   string                          `yaml:"filepath,omitempty"`
	Animations map[string]*animation.Animation `yaml:"animations"`

	// Ingame Related
	StateMachine    *state.StateMachine        `yaml:"-"`
	AnimationPlayer *animation.AnimationPlayer `yaml:"-"`
	AttackHasHit    bool                       `yaml:"-"` // prevents multi-hitting a single swing
}

type CharacterID int

const (
	Helmet CharacterID = iota
)

// LoadCharacter loads a character by its ID.
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
	c.AnimationPlayer = &animation.AnimationPlayer{}

	if c.Friction == 0 {
		panic(fmt.Sprintf("character %s has nil friction, check the character file", c.Name))
	}
	if c.JumpHeight == 0 {
		panic(fmt.Sprintf("character %s has nil jump height, check the character file", c.Name))
	}

	c.SetAnimation("idle")

	c.StateMachine = &state.StateMachine{}
	c.StateMachine.ActiveState = state.StateIdle
	c.StateMachine.HP = 10000
}

func (c *Character) SetAnimation(name string) {
	anim, ok := c.Animations[name]
	if !ok {
		if idleAnim, exists := c.Animations["idle"]; exists { // Fallback to idle animation
			anim = idleAnim
			c.AnimationPlayer.ShouldLoop = true
		} else {
			panic(fmt.Sprintf("Animation '%s' not found for character '%s' and no 'idle' fallback exists", name, c.Name))
		}

	}
	c.AnimationPlayer.ActiveAnimation = anim
	//c.AnimationPlayer.ShouldLoop = loop
	c.AnimationPlayer.FrameIndex = 0
	if len(anim.FrameData) > 0 {
		c.AnimationPlayer.FrameTimeLeft = anim.FrameData[0].Duration
	} else {
		c.AnimationPlayer.FrameTimeLeft = 0
	}
}

// ensureAnimation switches animations only when the requested one differs, preventing restarts.
func (c *Character) ensureAnimation(name string, loop bool) {
	current := ""
	if c.AnimationPlayer != nil && c.AnimationPlayer.ActiveAnimation != nil {
		current = c.AnimationPlayer.ActiveAnimation.Name
	}
	if current == name && c.AnimationPlayer.ShouldLoop == loop {
		return
	}

	c.SetAnimation(name)
	c.AnimationPlayer.ShouldLoop = loop
}

// updateAnimation advances the animation frame based on a simple frame counter
func (c *Character) updateAnimation() {
	if c.AnimationPlayer.ActiveAnimation == nil || len(c.AnimationPlayer.ActiveAnimation.FrameData) == 0 {
		return
	}

	c.AnimationPlayer.Update()

	// if anim ended and len(AnimationQueue) > 0, switch to next animation
	if c.AnimationPlayer.IsFinished() && len(c.AnimationPlayer.AnimationQueue) > 0 {
		c.SetAnimation(c.AnimationPlayer.AnimationQueue[0])
		c.AnimationPlayer.AnimationQueue = c.AnimationPlayer.AnimationQueue[1:]
	}
}

func loadCharacterByFile(filePath string) (*Character, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	character := &Character{}
	if err := yaml.Unmarshal(data, character); err != nil {
		return nil, err
	}

	// Convert relative paths to absolute paths based on YAML file location
	for _, anim := range character.Animations {
		for _, sprite := range anim.Sprites {
			if sprite.ImagePath != "" {
				sprite.ImagePath = resolveRelativePath(sprite.ImagePath, filePath)
			}
			ensureSpriteRect(sprite)
		}
	}
	return character, nil
}

// resolveRelativePath converts a relative path to an absolute path based on a reference path
func resolveRelativePath(relativePath, referencePath string) string {
	if filepath.IsAbs(relativePath) {
		return relativePath
	}
	referenceDir := filepath.Dir(referencePath)
	return filepath.Clean(filepath.Join(referenceDir, relativePath))
}

func (c *Character) Position() types.Vector2 {
	return c.StateMachine.Position
}

func (c *Character) Sprite() *animation.Sprite {
	if c == nil || c.AnimationPlayer == nil {
		return nil
	}
	return c.AnimationPlayer.ActiveSprite()
}

func (c *Character) BoundingBox() types.Rect {
	sprite := c.Sprite()
	return types.Rect{
		X: c.StateMachine.Position.X,
		Y: c.StateMachine.Position.Y,
		W: sprite.Rect.W,
		H: sprite.Rect.H,
	}
}

func ensureSpriteRect(sprite *animation.Sprite) {
	if sprite == nil {
		return
	}
	if sprite.Rect.W != 0 && sprite.Rect.H != 0 {
		return
	}

	width, height := probeImageSize(sprite.ImagePath)
	if width == 0 || height == 0 {
		return
	}

	sprite.Rect.W = float64(width)
	sprite.Rect.H = float64(height)
}

func probeImageSize(path string) (int, int) {
	if path == "" {
		return 0, 0
	}

	f, err := os.Open(path)
	if err != nil {
		return 0, 0
	}
	defer f.Close()

	conf, _, err := image.DecodeConfig(f)
	if err != nil {
		return 0, 0
	}

	return conf.Width, conf.Height
}
