package character

import (
	"bytes"
	"errors"
	"fgengine/animation"
	"fgengine/constants"
	"fgengine/types"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"gopkg.in/yaml.v3"
)

type Character struct {
	Name         string                  `yaml:"name" toml:"name"`
	Properties   CharacterProperties     `yaml:"properties" toml:"properties"`
	StateMachine *animation.StateMachine `yaml:"stateMachine" toml:"-"`
}

// CharacterProperties holds per-character tuning (SPEC §6.3). It is the
// only place character constants may live — never Go code (P2).
type CharacterProperties struct {
	MaxHP int `yaml:"maxHP,omitempty" toml:"maxHP,omitempty"`
}

func LoadCharacter(name string, playerSide int) (*Character, error) {
	chara, err := loadCharacterByName(name)
	if err != nil {
		return nil, err
	}
	chara.initialize(playerSide)
	return chara, nil
}

func loadCharacterByName(name string) (*Character, error) {
	filePath := "./assets/characters/" + name + ".yaml"
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read character file: %w", err)
	}

	character := &Character{
		Name: name,
	}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(character); err != nil {
		return nil, fmt.Errorf("failed to unmarshal character data: %w", err)
	}

	if character.StateMachine == nil || character.StateMachine.AnimPlayer == nil {
		return nil, fmt.Errorf("character file is missing stateMachine.activeAnim")
	}
	if character.StateMachine.AnimPlayer.Animations == nil {
		return nil, fmt.Errorf("character file is missing stateMachine.activeAnim.animations")
	}

	// Keep runtime animation names in sync with the map keys.
	for animName, anim := range character.StateMachine.AnimPlayer.Animations {
		if anim == nil {
			continue
		}
		anim.Name = animName

		for _, sprite := range anim.Sprites {
			if sprite.ImagePath != "" {
				sprite.ImagePath = resolveRelativePath(sprite.ImagePath, filePath)
			}
		}
	}
	// Sprite paths are resolved above, so Validate checks them on disk.
	if errs := character.Validate(); len(errs) > 0 {
		return nil, fmt.Errorf("invalid character file %s: %w", filePath, errors.Join(errs...))
	}
	return character, nil
}

func (c *Character) initialize(playerSide int) {
	if c.StateMachine == nil {
		c.StateMachine = new(animation.StateMachine{})
	}
	if c.StateMachine.AnimPlayer == nil {
		c.StateMachine.AnimPlayer = new(animation.AnimationPlayer{})
	}

	var initialX float64
	var facing animation.Orientation
	switch playerSide {
	case 1:
		initialX = constants.WorldWidth / 4
		facing = animation.Right
	case 2:
		initialX = 3 * constants.WorldWidth / 4
		facing = animation.Left
	}

	c.StateMachine.HP = c.Properties.MaxHP
	c.StateMachine.Position = types.Vector2{X: initialX, Y: constants.WorldHeight / 2}
	c.StateMachine.IsFacingLeft = facing
	c.StateMachine.Velocity = types.Vector2{}
	c.StateMachine.IgnoreGravityFrames = 0

	setInitialAnimation(c.StateMachine.AnimPlayer)

}

func setInitialAnimation(player *animation.AnimationPlayer) {
	if player == nil || len(player.Animations) == 0 || player.ActiveAnimation != nil {
		return
	}

	if _, ok := player.Animations["idle"]; ok {
		player.SetAnimation("idle")
		return
	}

	animNames := make([]string, 0, len(player.Animations))
	for name := range player.Animations {
		animNames = append(animNames, name)
	}
	slices.Sort(animNames)
	player.SetAnimation(animNames[0])
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
	return c.StateMachine.AnimPlayer.ActiveSprite()
}
