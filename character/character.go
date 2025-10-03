package character

import (
	"fgengine/animation"
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
	ID         int                             `yaml:"id,omitempty"`
	Name       string                          `yaml:"name"`
	Friction   float64                         `yaml:"friction,omitempty"`
	JumpHeight float64                         `yaml:"jumpHeight,omitempty"`
	FilePath   string                          `yaml:"filepath,omitempty"`
	Animations map[string]*animation.Animation `yaml:"animations"`

	// Ingame Related
	StateMachine *state.StateMachine `yaml:"-"`

	// animation related, TODO, move to a separate struct? maybe AnimationPlayer
	ActiveAnimation *animation.Animation `yaml:"-"`
	ActiveSprite    *animation.Sprite    `yaml:"-"`
	FrameIndex      int                  `yaml:"-"`
	SpriteIndex     int                  `yaml:"-"`
	ShouldLoop      bool                 `yaml:"-"`
	AnimationQueue  []string             `yaml:"-"`
	AnimationIndex  int                  `yaml:"-"`

	// frame counter
	FrameCounter int `yaml:"-"`
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
		// Try to fall back to a default animation if not found
		if idleAnim, exists := c.Animations["idle"]; exists {
			anim = idleAnim
			//fmt.Printf("Animation '%s' not found for character '%s', using 'idle' animation\n", name, c.Name)
		} else {
			fmt.Printf("Animation '%s' not found for character '%s' and no fallback available\n", name, c.Name)
			return
		}

	}

	// Only reset frame if switching to a different animation
	if c.ActiveAnimation != anim {
		c.ActiveAnimation = anim
		c.FrameIndex = 0
		c.FrameCounter = 0

		// Update sprite to match the first frame's SpriteIndex
		if len(anim.FrameData) > 0 {
			c.SpriteIndex = anim.FrameData[0].SpriteIndex
			if c.SpriteIndex >= 0 && c.SpriteIndex < len(anim.Sprites) {
				c.ActiveSprite = anim.Sprites[c.SpriteIndex]
			}
		} else if len(anim.Sprites) > 0 {
			// Fallback if no FrameData exists
			c.SpriteIndex = 0
			c.ActiveSprite = anim.Sprites[0]
		}
	}

	c.ShouldLoop = true
}

// updateAnimation advances the animation frame based on a simple frame counter
func (c *Character) updateAnimation() {
	if c.ActiveAnimation == nil || len(c.ActiveAnimation.FrameData) == 0 {
		return
	}

	c.FrameCounter++

	// Check if current frame duration has elapsed (using frame counter instead of time)
	currentFrameProps := &c.ActiveAnimation.FrameData[c.FrameIndex]
	if c.FrameCounter >= currentFrameProps.Duration {
		c.FrameCounter = 0

		// Check if current frame has an animation switch
		if currentFrameProps.AnimationSwitch != "" {
			// Switch to the specified animation
			if newAnimation, exists := c.Animations[currentFrameProps.AnimationSwitch]; exists {
				c.ActiveAnimation = newAnimation
				c.FrameIndex = 0
				c.FrameCounter = 0
				// Update sprite to match the first frame's SpriteIndex
				if len(newAnimation.FrameData) > 0 {
					c.SpriteIndex = newAnimation.FrameData[0].SpriteIndex
					if c.SpriteIndex >= 0 && c.SpriteIndex < len(newAnimation.Sprites) {
						c.ActiveSprite = newAnimation.Sprites[c.SpriteIndex]
					}
				} else if len(newAnimation.Sprites) > 0 {
					c.SpriteIndex = 0
					c.ActiveSprite = newAnimation.Sprites[0]
				}
				return // Early return to avoid normal frame advancement
			} else {
				// Animation switch target not found, try fallback
				fmt.Printf("AnimationSwitch target '%s' not found for character '%s'\n", currentFrameProps.AnimationSwitch, c.Name)
				if fallbackAnim, exists := c.Animations["notFound"]; exists {
					c.ActiveAnimation = fallbackAnim
					c.FrameIndex = 0
					c.FrameCounter = 0
					if len(fallbackAnim.FrameData) > 0 {
						c.SpriteIndex = fallbackAnim.FrameData[0].SpriteIndex
						if c.SpriteIndex >= 0 && c.SpriteIndex < len(fallbackAnim.Sprites) {
							c.ActiveSprite = fallbackAnim.Sprites[c.SpriteIndex]
						}
					} else if len(fallbackAnim.Sprites) > 0 {
						c.SpriteIndex = 0
						c.ActiveSprite = fallbackAnim.Sprites[0]
					}
					return
				} else if idleAnim, exists := c.Animations["idle"]; exists {
					c.ActiveAnimation = idleAnim
					c.FrameIndex = 0
					c.FrameCounter = 0
					if len(idleAnim.FrameData) > 0 {
						c.SpriteIndex = idleAnim.FrameData[0].SpriteIndex
						if c.SpriteIndex >= 0 && c.SpriteIndex < len(idleAnim.Sprites) {
							c.ActiveSprite = idleAnim.Sprites[c.SpriteIndex]
						}
					} else if len(idleAnim.Sprites) > 0 {
						c.SpriteIndex = 0
						c.ActiveSprite = idleAnim.Sprites[0]
					}
					return
				}
			}
		}

		c.FrameIndex++

		// Handle end of animation
		if c.FrameIndex >= len(c.ActiveAnimation.FrameData) {
			if c.ShouldLoop {
				c.FrameIndex = 0
			} else {
				c.FrameIndex = len(c.ActiveAnimation.FrameData) - 1 // Stay on last frame
			}
		}

		// Update sprite based on current frame's SpriteIndex
		if c.FrameIndex < len(c.ActiveAnimation.FrameData) {
			c.SpriteIndex = c.ActiveAnimation.FrameData[c.FrameIndex].SpriteIndex
			if c.SpriteIndex >= 0 && c.SpriteIndex < len(c.ActiveAnimation.Sprites) {
				c.ActiveSprite = c.ActiveAnimation.Sprites[c.SpriteIndex]
			}
		}
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

// GetCurrentFrameProperties returns the frame properties for the current frame of the active animation
func (c *Character) GetCurrentFrameProperties() *animation.FrameData {
	if c.ActiveAnimation == nil || c.FrameIndex >= len(c.ActiveAnimation.FrameData) {
		return nil
	}
	return &c.ActiveAnimation.FrameData[c.FrameIndex]
}
